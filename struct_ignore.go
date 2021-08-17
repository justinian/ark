package main

func readIgnoredStruct(dataSize int, vr valueReader) (Property, error) {
	vr.skip(dataSize)
	return nil, nil
}

func init() {
	addStructType("ItemNetID", readIgnoredStruct)
}
