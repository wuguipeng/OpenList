package _139

import (
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
)

type Addition struct {
	//Account       string `json:"account" required:"true"`
	Authorization string `json:"authorization" type:"text" help:"Authorization can be used alone. If empty, use mail_cookies alone for fast login, or mail_cookies + username + password for full login fallback."`
	Username      string `json:"username" help:"Use together with password and mail_cookies for full login fallback when Authorization is empty."`
	Password      string `json:"password" secret:"true" help:"Use together with username and mail_cookies for full login fallback when Authorization is empty."`
	SmsCode       string `json:"sms_code" secret:"true" help:"Fill this only after OpenList reports that a 139 Mail SMS verification code was sent, then save the storage again."`
	MailCookies   string `json:"mail_cookies" type:"text" help:"Optional cookies from mail.10086.cn. Used alone for fast login, or together with username and password for full login fallback."`
	driver.RootID
	Type                 string `json:"type" type:"select" options:"personal_new,family,group,personal,share" default:"personal_new"`
	LinkID               string `json:"link_id" type:"text" help:"Multiple shares are separated by commas or new lines. Use link_id#password for password-protected shares."`
	CloudID              string `json:"cloud_id"`
	UserDomainID         string `json:"user_domain_id" help:"ud_id in Cookie, fill in to show disk usage"`
	CustomUploadPartSize int64  `json:"custom_upload_part_size" type:"number" default:"0" help:"0 for auto"`
	ReportRealSize       bool   `json:"report_real_size" type:"bool" default:"true" help:"Enable to report the real file size during upload"`
	UseLargeThumbnail    bool   `json:"use_large_thumbnail" type:"bool" default:"false" help:"Enable to use large thumbnail for images"`
	UseOldStreamUpload   bool   `json:"use_old_stream_upload" type:"bool" default:"false" help:"Enable to use old stream upload method (not support rapid upload)"`
	ClientType           string `json:"client_type" type:"select" options:"web,pc" default:"web" help:"Client type reported to the personal cloud API. web (9) allows uploading files up to 50 GB; pc (13) raises the single file upload limit to 500 GB."`
}

var config = driver.Config{
	Name:             "139Yun",
	LocalSort:        true,
	ProxyRangeOption: true,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		d := &Yun139{}
		d.ProxyRange = true
		return d
	})
}
