package modules

/*
Usage
- err := NewNotification(&Email{}).Push()
- err := NewNotification(&Mobile{}).Push()
*/

import "fmt"

type INotification interface {
	Push() error
}

type Notification struct {
	email_info  *Email
	mobile_info *Mobile
}

func NewNotification(ntype interface{}) INotification {
	n := &Notification{}
	if e, ok := ntype.(*Email); ok {
		n.emailSet(e)
	} else if m, ok := ntype.(*Mobile); ok {
		n.mobileSet(m)
	} else {
		panic("type issue")
	}
	return n
}

func (n *Notification) emailSet(paras *Email) {
	n.email_info = paras
}

func (n *Notification) mobileSet(paras *Mobile) {
	n.mobile_info = paras
}

func (n *Notification) Push() error {
	if n.email_info != nil {
		fmt.Println("push email")
	} else if n.mobile_info != nil {
		fmt.Println("push mobile")
	}
	return nil
}
