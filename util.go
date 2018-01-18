package main

func incamount(scale int) int {
	    if (ACCESS_FLAG(F_DF)) {/* down */
		    return scale * -1
	    }
	return scale * 1
}


