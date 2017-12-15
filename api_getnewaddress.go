package multichain

func (client *Client) GetNewAddress() (Response, error) {

	msg := client.ChainMsg(
		"getnewaddress",
		[]interface{}{},
	)

	obj, err := client.post(msg)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
