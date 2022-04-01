package service

//初始化CountMap
func SetCCTVCountMap(cctvCountMap map[string]int, cctvNo string) {

	if _, ok := cctvCountMap[cctvNo]; !ok {
		cctvCountMap[cctvNo] = 1
	}
}
