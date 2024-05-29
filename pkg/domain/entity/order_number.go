package entity

type OrderNumber string

func (b *OrderNumber) String() string {
	if b == nil {
		return ""
	}
	return string(*b)
}

func (b *OrderNumber) FromBytes(d []byte) (*OrderNumber, error) {
	s := OrderNumber(d)
	return &s, nil
}

//func (b *OrderNumber) FromBytes(s []byte) (*OrderNumber, error) {
//	t := new(OrderNumber)
//	err := t.Scan(s)
//	return t, err
//}
//func (b *OrderNumber) TextValue() (pgtype.Text, error) {
//	if b != nil {
//		s := (*big.Int)(b).String()
//		return pgtype.Text{String: s, Valid: true}, nil
//	}
//	return pgtype.Text{}, nil
//}
//
//func (b *OrderNumber) Value() (driver.Value, error) {
//	if b != nil {
//		s := (*big.Int)(b).String()
//		return driver.Value(s), nil
//	}
//	return pgtype.Text{}, nil
//}
//
//func (b *OrderNumber) String() string {
//	if b != nil {
//		return (*big.Int)(b).String()
//	}
//	return ""
//}
//
//func (b *OrderNumber) ScanText(value pgtype.Text) error {
//	return b.Scan(value.String)
//}
//
//func (b *OrderNumber) Scan(value interface{}) error {
//	if value == nil {
//		b = nil
//	}
//
//	switch t := value.(type) {
//	case []uint8:
//		_, ok := (*big.Int)(b).SetString(string(value.([]uint8)), 10)
//		if !ok {
//			return fmt.Errorf("failed to load value to []uint8: %v", value)
//		}
//	case string:
//		_, ok := (*big.Int)(b).SetString(value.(string), 10)
//		if !ok {
//			return fmt.Errorf("failed to load value to string: %v", value)
//		}
//	default:
//		return fmt.Errorf("could not scan type %T into OrderNumber", t)
//	}
//
//	return nil
//}
