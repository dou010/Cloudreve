package model

import (
	"github.com/HFO4/cloudreve/pkg/conf"
	"github.com/HFO4/cloudreve/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mcuadros/go-version"
	"io/ioutil"
)

//执行数据迁移
func migration() {
	// 检查 version.lock 确认是否需要执行迁移
	// Debug 模式及测试模式下一定会执行迁移
	if !conf.SystemConfig.Debug && gin.Mode() != gin.TestMode {
		if util.Exists("version.lock") {
			versionLock, _ := ioutil.ReadFile("version.lock")
			if version.Compare(string(versionLock), conf.BackendVersion, "=") {
				util.Log().Info("后端版本匹配，跳过数据库迁移")
				return
			}
		}
	}

	util.Log().Info("开始进行数据库自动迁移...")

	// 自动迁移模式
	if conf.DatabaseConfig.Type == "mysql" {
		DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	DB.AutoMigrate(&User{}, &Setting{}, &Group{}, &Policy{}, &Folder{}, &File{}, &StoragePack{}, &Share{},
		&Task{}, &Download{}, &Tag{}, &Webdav{}, &Order{}, &Redeem{})

	// 创建初始存储策略
	addDefaultPolicy()

	// 创建初始用户组
	addDefaultGroups()

	// 创建初始管理员账户
	addDefaultUser()

	// 向设置数据表添加初始设置
	addDefaultSettings()

	// 迁移完毕后写入版本锁 version.lock
	err := conf.WriteVersionLock()
	if err != nil {
		util.Log().Warning("无法写入版本控制锁 version.lock, %s", err)
	}

	util.Log().Info("数据库自动迁移结束")

}

func addDefaultPolicy() {
	_, err := GetPolicyByID(uint(1))
	// 未找到初始存储策略时，则创建
	if gorm.IsRecordNotFoundError(err) {
		defaultPolicy := Policy{
			Name:               "默认存储策略",
			Type:               "local",
			Server:             "/api/v3/file/upload",
			BaseURL:            "http://cloudreve.org/public/uploads/",
			MaxSize:            10 * 1024 * 1024 * 1024,
			AutoRename:         true,
			DirNameRule:        "uploads/{uid}/{path}",
			FileNameRule:       "{uid}_{randomkey8}_{originname}",
			IsOriginLinkEnable: false,
			OptionsSerialized: PolicyOption{
				FileType: []string{},
			},
		}
		if err := DB.Create(&defaultPolicy).Error; err != nil {
			util.Log().Panic("无法创建初始存储策略, %s", err)
		}
	}
}

func addDefaultSettings() {
	defaultSettings := []Setting{
		{Name: "siteURL", Value: ``, Type: "basic"},
		{Name: "siteName", Value: `Cloudreve`, Type: "basic"},
		{Name: "siteStatus", Value: `open`, Type: "basic"},
		{Name: "register_enabled", Value: `1`, Type: "register"},
		{Name: "default_group", Value: `2`, Type: "register"},
		{Name: "siteKeywords", Value: `网盘，网盘`, Type: "basic"},
		{Name: "siteDes", Value: `Cloudreve`, Type: "basic"},
		{Name: "siteTitle", Value: `平步云端`, Type: "basic"},
		{Name: "fromName", Value: `Cloudreve`, Type: "mail"},
		{Name: "mail_keepalive", Value: `30`, Type: "mail"},
		{Name: "fromAdress", Value: `no-reply@acg.blue`, Type: "mail"},
		{Name: "smtpHost", Value: `smtp.mxhichina.com`, Type: "mail"},
		{Name: "smtpPort", Value: `25`, Type: "mail"},
		{Name: "replyTo", Value: `abslant@126.com`, Type: "mail"},
		{Name: "smtpUser", Value: `no-reply@acg.blue`, Type: "mail"},
		{Name: "smtpPass", Value: ``, Type: "mail"},
		{Name: "encriptionType", Value: `no`, Type: "mail"},
		{Name: "over_used_template", Value: `<meta name="viewport"content="width=device-width"><meta http-equiv="Content-Type"content="text/html; charset=UTF-8"><title>容量超额提醒</title><style type="text/css">img{max-width:100%}body{-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:none;width:100%!important;height:100%;line-height:1.6em}body{background-color:#f6f6f6}@media only screen and(max-width:640px){body{padding:0!important}h1{font-weight:800!important;margin:20px 0 5px!important}h2{font-weight:800!important;margin:20px 0 5px!important}h3{font-weight:800!important;margin:20px 0 5px!important}h4{font-weight:800!important;margin:20px 0 5px!important}h1{font-size:22px!important}h2{font-size:18px!important}h3{font-size:16px!important}.container{padding:0!important;width:100%!important}.content{padding:0!important}.content-wrap{padding:10px!important}.invoice{width:100%!important}}</style><table class="body-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; background-color: #f6f6f6; margin: 0;"bgcolor="#f6f6f6"><tbody><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td><td class="container"width="600"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; display: block !important; max-width: 600px !important; clear: both !important; margin: 0 auto;"valign="top"><div class="content"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; max-width: 600px; display: block; margin: 0 auto; padding: 20px;"><table class="main"width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; border-radius: 3px; background-color: #fff; margin: 0; border: 1px 
solid #e9e9e9;"bgcolor="#fff"><tbody><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="alert alert-warning"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 16px; vertical-align: top; color: #fff; font-weight: 500; text-align: center; border-radius: 3px 3px 0 0; background-color: #FF9F00; margin: 0; padding: 20px;"align="center"bgcolor="#FF9F00"valign="top">容量超额警告</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 20px;"valign="top"><table width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tbody><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">亲爱的<strong style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">{userName}</strong>：</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">由于{notifyReason}，您在{siteTitle}的账户的容量使用超出配额，您将无法继续上传新文件，请尽快清理文件，否则我们将会禁用您的账户。</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top"><a href="{siteUrl}Login"class="btn-primary"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; color: #FFF; text-decoration: none; line-height: 2em; font-weight: bold; text-align: center; cursor: pointer; display: inline-block; border-radius: 5px; text-transform: capitalize; background-color: #348eda; margin: 0; border-color: #348eda; border-style: solid; border-width: 10px 20px;">登录{siteTitle}</a></td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">感谢您选择{siteTitle}。</td></tr></tbody></table></td></tr></tbody></table><div class="footer"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; clear: both; color: #999; margin: 0; padding: 20px;"><table width="100%"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tbody><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="aligncenter content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 12px; vertical-align: top; color: #999; text-align: center; margin: 0; padding: 0 0 20px;"align="center"valign="top">此邮件由系统自动发送，请不要直接回复。</td></tr></tbody></table></div></div></td><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td></tr></tbody></table>`, Type: "mail_template"},
		{Name: "ban_time", Value: `10`, Type: "storage_policy"},
		{Name: "maxEditSize", Value: `100000`, Type: "file_edit"},
		{Name: "oss_timeout", Value: `3600`, Type: "timeout"},
		{Name: "archive_timeout", Value: `30`, Type: "timeout"},
		{Name: "download_timeout", Value: `30`, Type: "timeout"},
		{Name: "preview_timeout", Value: `60`, Type: "timeout"},
		{Name: "doc_preview_timeout", Value: `60`, Type: "timeout"},
		{Name: "upload_credential_timeout", Value: `1800`, Type: "timeout"},
		{Name: "upload_session_timeout", Value: `86400`, Type: "timeout"},
		{Name: "slave_api_timeout", Value: `60`, Type: "timeout"},
		{Name: "onedrive_monitor_timeout", Value: `600`, Type: "timeout"},
		{Name: "share_download_session_timeout", Value: `2073600`, Type: "timeout"},
		{Name: "onedrive_callback_check", Value: `20`, Type: "timeout"},
		{Name: "aria2_call_timeout", Value: `5`, Type: "timeout"},
		{Name: "onedrive_chunk_retries", Value: `1`, Type: "retry"},
		{Name: "allowdVisitorDownload", Value: `false`, Type: "share"},
		{Name: "login_captcha", Value: `0`, Type: "login"},
		{Name: "qq_login", Value: `0`, Type: "login"},
		{Name: "qq_login_id", Value: ``, Type: "login"},
		{Name: "qq_login_key", Value: ``, Type: "login"},
		{Name: "reg_captcha", Value: `0`, Type: "login"},
		{Name: "email_active", Value: `0`, Type: "register"},
		{Name: "mail_activation_template", Value: `<!DOCTYPE html PUBLIC"-//W3C//DTD XHTML 1.0 Transitional//EN""http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><html xmlns="http://www.w3.org/1999/xhtml"style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; box-sizing: border-box; 
font-size: 14px; margin: 0;"><head><meta name="viewport"content="width=device-width"/><meta http-equiv="Content-Type"content="text/html; charset=UTF-8"/><title>容量超额提醒</title><style type="text/css">img{max-width:100%}body{-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:none;width:100%!important;height:100%;line-height:1.6em}body{background-color:#f6f6f6}@media only screen and(max-width:640px){body{padding:0!important}h1{font-weight:800!important;margin:20px 0 5px!important}h2{font-weight:800!important;margin:20px 0 5px!important}h3{font-weight:800!important;margin:20px 0 5px!important}h4{font-weight:800!important;margin:20px 0 5px!important}h1{font-size:22px!important}h2{font-size:18px!important}h3{font-size:16px!important}.container{padding:0!important;width:100%!important}.content{padding:0!important}.content-wrap{padding:10px!important}.invoice{width:100%!important}}</style></head><body itemscope itemtype="http://schema.org/EmailMessage"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: 
border-box; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; width: 100% !important; height: 100%; line-height: 1.6em; background-color: #f6f6f6; margin: 0;"bgcolor="#f6f6f6"><table class="body-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; background-color: #f6f6f6; margin: 0;"bgcolor="#f6f6f6"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; 
box-sizing: border-box; font-size: 14px; margin: 0;"><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td><td class="container"width="600"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; display: block !important; max-width: 600px !important; clear: both !important; margin: 0 auto;"valign="top"><div class="content"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; max-width: 600px; display: block; margin: 0 auto; padding: 20px;"><table class="main"width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; border-radius: 3px; background-color: #fff; margin: 0; border: 1px 
solid #e9e9e9;"bgcolor="#fff"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 
14px; margin: 0;"><td class="alert alert-warning"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 16px; vertical-align: top; color: #fff; font-weight: 500; text-align: center; border-radius: 3px 3px 0 0; background-color: #009688; margin: 0; padding: 20px;"align="center"bgcolor="#FF9F00"valign="top">激活{siteTitle}账户</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 20px;"valign="top"><table width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica 
Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">亲爱的<strong style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">{userName}</strong>：</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">感谢您注册{siteTitle},请点击下方按钮完成账户激活。</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top"><a href="{activationUrl}"class="btn-primary"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; color: #FFF; text-decoration: none; line-height: 2em; font-weight: bold; text-align: center; cursor: pointer; display: inline-block; border-radius: 5px; text-transform: capitalize; background-color: #009688; margin: 0; border-color: #009688; border-style: solid; border-width: 10px 20px;">激活账户</a></td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">感谢您选择{siteTitle}。</td></tr></table></td></tr></table><div class="footer"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; clear: both; color: #999; margin: 0; padding: 20px;"><table width="100%"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="aligncenter content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 12px; vertical-align: top; color: #999; text-align: center; margin: 0; padding: 0 0 20px;"align="center"valign="top">此邮件由系统自动发送，请不要直接回复。</td></tr></table></div></div></td><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td></tr></table></body></html>`, Type: "mail_template"},
		{Name: "forget_captcha", Value: `0`, Type: "login"},
		{Name: "mail_reset_pwd_template", Value: `<!DOCTYPE html PUBLIC"-//W3C//DTD XHTML 1.0 Transitional//EN""http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><html xmlns="http://www.w3.org/1999/xhtml"style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; box-sizing: border-box; 
font-size: 14px; margin: 0;"><head><meta name="viewport"content="width=device-width"/><meta http-equiv="Content-Type"content="text/html; charset=UTF-8"/><title>重设密码</title><style type="text/css">img{max-width:100%}body{-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:none;width:100%!important;height:100%;line-height:1.6em}body{background-color:#f6f6f6}@media only screen and(max-width:640px){body{padding:0!important}h1{font-weight:800!important;margin:20px 0 5px!important}h2{font-weight:800!important;margin:20px 0 5px!important}h3{font-weight:800!important;margin:20px 0 5px!important}h4{font-weight:800!important;margin:20px 0 5px!important}h1{font-size:22px!important}h2{font-size:18px!important}h3{font-size:16px!important}.container{padding:0!important;width:100%!important}.content{padding:0!important}.content-wrap{padding:10px!important}.invoice{width:100%!important}}</style></head><body itemscope itemtype="http://schema.org/EmailMessage"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: 
border-box; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; width: 100% !important; height: 100%; line-height: 1.6em; background-color: #f6f6f6; margin: 0;"bgcolor="#f6f6f6"><table class="body-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; background-color: #f6f6f6; margin: 0;"bgcolor="#f6f6f6"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; 
box-sizing: border-box; font-size: 14px; margin: 0;"><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td><td class="container"width="600"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; display: block !important; max-width: 600px !important; clear: both !important; margin: 0 auto;"valign="top"><div class="content"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; max-width: 600px; display: block; margin: 0 auto; padding: 20px;"><table class="main"width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; border-radius: 3px; background-color: #fff; margin: 0; border: 1px 
solid #e9e9e9;"bgcolor="#fff"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 
14px; margin: 0;"><td class="alert alert-warning"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 16px; vertical-align: top; color: #fff; font-weight: 500; text-align: center; border-radius: 3px 3px 0 0; background-color: #2196F3; margin: 0; padding: 20px;"align="center"bgcolor="#FF9F00"valign="top">重设{siteTitle}密码</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-wrap"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 20px;"valign="top"><table width="100%"cellpadding="0"cellspacing="0"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica 
Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">亲爱的<strong style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">{userName}</strong>：</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">请点击下方按钮完成密码重设。如果非你本人操作，请忽略此邮件。</td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top"><a href="{resetUrl}"class="btn-primary"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; color: #FFF; text-decoration: none; line-height: 2em; font-weight: bold; text-align: center; cursor: pointer; display: inline-block; border-radius: 5px; text-transform: capitalize; background-color: #2196F3; margin: 0; border-color: #2196F3; border-style: solid; border-width: 10px 20px;">重设密码</a></td></tr><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;"valign="top">感谢您选择{siteTitle}。</td></tr></table></td></tr></table><div class="footer"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; clear: both; color: #999; margin: 0; padding: 20px;"><table width="100%"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;"><td class="aligncenter content-block"style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 12px; vertical-align: top; color: #999; text-align: center; margin: 0; padding: 0 0 20px;"align="center"valign="top">此邮件由系统自动发送，请不要直接回复。</td></tr></table></div></div></td><td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;"valign="top"></td></tr></table></body></html>`, Type: "mail_template"},
		{Name: "pack_data", Value: `[]`, Type: "pack"},
		{Name: "database_version", Value: `6`, Type: "version"},
		{Name: "alipay_enabled", Value: `0`, Type: "payment"},
		{Name: "payjs_enabled", Value: `0`, Type: "payment"},
		{Name: "payjs_id", Value: ``, Type: "payment"},
		{Name: "payjs_secret", Value: ``, Type: "payment"},
		{Name: "appid", Value: ``, Type: "payment"},
		{Name: "appkey", Value: ``, Type: "payment"},
		{Name: "shopid", Value: ``, Type: "payment"},
		{Name: "hot_share_num", Value: `10`, Type: "share"},
		{Name: "allow_buy_group", Value: `1`, Type: "group_sell"},
		{Name: "group_sell_data", Value: `[]`, Type: "group_sell"},
		{Name: "gravatar_server", Value: `https://gravatar.loli.net/`, Type: "avatar"},
		{Name: "defaultTheme", Value: `#3f51b5`, Type: "basic"},
		{Name: "themes", Value: `{"#3f51b5":{"palette":{"primary":{"main":"#3f51b5"},"secondary":{"main":"#f50057"}}},"#2196f3":{"palette":{"primary":{"main":"#2196f3"},"secondary":{"main":"#FFC107"}}},"#673AB7":{"palette":{"primary":{"main":"#673AB7"},"secondary":{"main":"#2196F3"}}},"#E91E63":{"palette":{"primary":{"main":"#E91E63"},"secondary":{"main":"#42A5F5","contrastText":"#fff"}}},"#FF5722":{"palette":{"primary":{"main":"#FF5722"},"secondary":{"main":"#3F51B5"}}},"#FFC107":{"palette":{"primary":{"main":"#FFC107"},"secondary":{"main":"#26C6DA"}}},"#8BC34A":{"palette":{"primary":{"main":"#8BC34A","contrastText":"#fff"},"secondary":{"main":"#FF8A65","contrastText":"#fff"}}},"#009688":{"palette":{"primary":{"main":"#009688"},"secondary":{"main":"#4DD0E1","contrastText":"#fff"}}},"#607D8B":{"palette":{"primary":{"main":"#607D8B"},"secondary":{"main":"#F06292"}}},"#795548":{"palette":{"primary":{"main":"#795548"},"secondary":{"main":"#4CAF50","contrastText":"#fff"}}}}`, Type: "basic"},
		{Name: "aria2_token", Value: `your token`, Type: "aria2"},
		{Name: "aria2_token", Value: `your token`, Type: "aria2"},
		{Name: "aria2_temp_path", Value: ``, Type: "aria2"},
		{Name: "aria2_options", Value: `[]`, Type: "aria2"},
		{Name: "aria2_interval", Value: `10`, Type: "aria2"},
		{Name: "max_worker_num", Value: `10`, Type: "task"},
		{Name: "max_parallel_transfer", Value: `4`, Type: "task"},
		{Name: "secret_key", Value: util.RandStringRunes(256), Type: "auth"},
		{Name: "temp_path", Value: "temp", Type: "path"},
		{Name: "avatar_path", Value: "avatar", Type: "path"},
		{Name: "avatar_size", Value: "2097152", Type: "avatar"},
		{Name: "avatar_size_l", Value: "200", Type: "avatar"},
		{Name: "avatar_size_m", Value: "130", Type: "avatar"},
		{Name: "avatar_size_s", Value: "50", Type: "avatar"},
		{Name: "score_enabled", Value: "1", Type: "score"},
		{Name: "share_score_rate", Value: "80", Type: "score"},
		{Name: "score_price", Value: "1", Type: "score"},
		{Name: "home_view_method", Value: "icon", Type: "view"},
		{Name: "share_view_method", Value: "list", Type: "view"},
		{Name: "cron_garbage_collect", Value: "@hourly", Type: "cron"},
		{Name: "cron_notify_user", Value: "@hourly", Type: "cron"},
		{Name: "cron_ban_user", Value: "@hourly", Type: "cron"},
		{Name: "authn_enabled", Value: "1", Type: "authn"},
	}

	for _, value := range defaultSettings {
		DB.Where(Setting{Name: value.Name}).Create(&value)
	}
}

func addDefaultGroups() {
	_, err := GetGroupByID(1)
	// 未找到初始管理组时，则创建
	if gorm.IsRecordNotFoundError(err) {
		defaultAdminGroup := Group{
			Name:          "管理员",
			PolicyList:    []uint{1},
			MaxStorage:    1 * 1024 * 1024 * 1024,
			ShareEnabled:  true,
			Color:         "danger",
			WebDAVEnabled: true,
			OptionsSerialized: GroupOption{
				ArchiveDownload: true,
				ArchiveTask:     true,
				ShareDownload:   true,
			},
		}
		if err := DB.Create(&defaultAdminGroup).Error; err != nil {
			util.Log().Panic("无法创建管理用户组, %s", err)
		}
	}

	err = nil
	_, err = GetGroupByID(2)
	// 未找到初始注册会员时，则创建
	if gorm.IsRecordNotFoundError(err) {
		defaultAdminGroup := Group{
			Name:          "注册会员",
			PolicyList:    []uint{1},
			MaxStorage:    1 * 1024 * 1024 * 1024,
			ShareEnabled:  true,
			Color:         "danger",
			WebDAVEnabled: true,
		}
		if err := DB.Create(&defaultAdminGroup).Error; err != nil {
			util.Log().Panic("无法创建初始注册会员用户组, %s", err)
		}
	}

	err = nil
	_, err = GetGroupByID(3)
	// 未找到初始游客用户组时，则创建
	if gorm.IsRecordNotFoundError(err) {
		defaultAdminGroup := Group{
			Name:     "游客",
			Policies: "[]",
		}
		if err := DB.Create(&defaultAdminGroup).Error; err != nil {
			util.Log().Panic("无法创建初始游客用户组, %s", err)
		}
	}
}

func addDefaultUser() {
	_, err := GetUserByID(1)

	// 未找到初始用户时，则创建
	if gorm.IsRecordNotFoundError(err) {
		defaultUser := NewUser()
		//TODO 动态生成密码
		defaultUser.Email = "admin@cloudreve.org"
		defaultUser.Nick = "admin"
		defaultUser.Status = Active
		defaultUser.GroupID = 1
		err := defaultUser.SetPassword("admin")
		if err != nil {
			util.Log().Panic("无法创建密码, %s", err)
		}
		if err := DB.Create(&defaultUser).Error; err != nil {
			util.Log().Panic("无法创建初始用户, %s", err)
		}
	}
}
