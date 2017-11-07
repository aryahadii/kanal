package handler

func removeItemFromStringSlice(item string, slice *[]string) bool {
	for i, sliceItem := range *slice {
		if sliceItem == item {
			(*slice)[i] = (*slice)[len(*slice)-1]
			*slice = (*slice)[:len(*slice)-1]
			return true
		}
	}
	return false
}
