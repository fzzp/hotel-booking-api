package rdb

func smsCodeKey(phone string) string {
	return "sms:" + phone
}
