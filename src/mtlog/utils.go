package mtlog

import "time"

func getLogTime() []byte {
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	nanosecond := now.Nanosecond() / 1000

	var buf [32]byte
	bp := len(buf) - 1

	wid := 6
	for nanosecond >= 10 || wid > 1 {
		wid--
		q := nanosecond / 10
		buf[bp] = byte('0' + nanosecond - q*10)
		bp--
		nanosecond = q
	}
	buf[bp] = byte('0' + nanosecond)
	bp--
	buf[bp] = byte('.')
	bp--

	wid = 2
	for second >= 10 || wid > 1 {
		wid--
		q := second / 10
		buf[bp] = byte('0' + second - q*10)
		bp--
		second = q
	}
	buf[bp] = byte('0' + second)
	bp--
	buf[bp] = byte(':')
	bp--

	wid = 2
	for minute >= 10 || wid > 1 {
		wid--
		q := minute / 10
		buf[bp] = byte('0' + minute - q*10)
		bp--
		minute = q
	}
	buf[bp] = byte('0' + minute)
	bp--
	buf[bp] = byte(':')
	bp--

	wid = 2
	for hour >= 10 || wid > 1 {
		wid--
		q := hour / 10
		buf[bp] = byte('0' + hour - q*10)
		bp--
		hour = q
	}
	buf[bp] = byte('0' + hour)
	bp--
	buf[bp] = byte(' ')
	bp--

	wid = 2
	for day >= 10 || wid > 1 {
		wid--
		q := day / 10
		buf[bp] = byte('0' + day - q*10)
		bp--
		day = q
	}
	buf[bp] = byte('0' + day)
	bp--
	buf[bp] = byte('-')
	bp--

	wid = 2
	for month >= 10 || wid > 1 {
		wid--
		q := month / 10
		buf[bp] = byte('0' + month - q*10)
		bp--
		month = q
	}
	buf[bp] = byte('0' + month)
	bp--
	buf[bp] = byte('-')
	bp--

	wid = 4
	for year >= 10 || wid > 1 {
		wid--
		q := year / 10
		buf[bp] = byte('0' + year - q*10)
		bp--
		year = q
	}
	buf[bp] = byte('0' + year)

	return buf[bp:]
}

func getFileTime() []byte {
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	nanosecond := now.Nanosecond() / 1000

	var buf [32]byte
	bp := len(buf) - 1

	wid := 6
	for nanosecond >= 10 || wid > 1 {
		wid--
		q := nanosecond / 10
		buf[bp] = byte('0' + nanosecond - q*10)
		bp--
		nanosecond = q
	}
	buf[bp] = byte('0' + nanosecond)
	bp--

	wid = 2
	for second >= 10 || wid > 1 {
		wid--
		q := second / 10
		buf[bp] = byte('0' + second - q*10)
		bp--
		second = q
	}
	buf[bp] = byte('0' + second)
	bp--

	wid = 2
	for minute >= 10 || wid > 1 {
		wid--
		q := minute / 10
		buf[bp] = byte('0' + minute - q*10)
		bp--
		minute = q
	}
	buf[bp] = byte('0' + minute)
	bp--

	wid = 2
	for hour >= 10 || wid > 1 {
		wid--
		q := hour / 10
		buf[bp] = byte('0' + hour - q*10)
		bp--
		hour = q
	}
	buf[bp] = byte('0' + hour)
	bp--

	wid = 2
	for day >= 10 || wid > 1 {
		wid--
		q := day / 10
		buf[bp] = byte('0' + day - q*10)
		bp--
		day = q
	}
	buf[bp] = byte('0' + day)
	bp--

	wid = 2
	for month >= 10 || wid > 1 {
		wid--
		q := month / 10
		buf[bp] = byte('0' + month - q*10)
		bp--
		month = q
	}
	buf[bp] = byte('0' + month)
	bp--

	wid = 4
	for year >= 10 || wid > 1 {
		wid--
		q := year / 10
		buf[bp] = byte('0' + year - q*10)
		bp--
		year = q
	}
	buf[bp] = byte('0' + year)

	return buf[bp:]
}
