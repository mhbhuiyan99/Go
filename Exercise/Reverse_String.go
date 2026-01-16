// Problem: https://app.gointerview.dev/challenge/2

func ReverseString(s string) (res_s string) {
	
	// solution - 1
	for i := len(s) -1; i >= 0; i-- {
	    res_s = res_s + string(s[i])
	}
	
	// solution - 2
	slice := make([]string, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
	    slice = append(slice,string(s[i]))
	}
	res_s = strings.Join(slice, "") 
	
	// solution - 3
	r := []rune(s) 
	for i, j := 0, len(r) - 1; i < j; i, j = i+1, j-1 {
	    r[i], r[j] = r[j], r[i]
	}
	res_s = string(r)
	
	return
}
