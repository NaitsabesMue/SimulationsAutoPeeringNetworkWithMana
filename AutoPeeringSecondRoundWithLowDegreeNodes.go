package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const totalNodes uint64 = 10000   // # of nodes
const zipfs float64 = 0.8         // zipfs parameter
const k uint32 = 4                // Number of in/out neighbors
const inBound uint32 = k          // InNeighbours
const outBound uint32 = k         // OutNeighbours
const loopNumber1 uint64 = 150000 // Loop First Round
const loopNumber2 uint64 = 100000 // Loop Second Round
const bigR uint64 = 20            // (j-i)<= bigR
const rho uint64 = 2              // mana[i] < mana[j]*rho

func main() {
	var zipfsMana = make(map[uint64]uint64, totalNodes)
	var listPossibleNeighbors = make(map[uint64][]uint64, totalNodes)
	createZipfsMana(zipfsMana)
	inNeighbors := make(map[uint64][]uint64, totalNodes)
	outNeighbors := make(map[uint64][]uint64, totalNodes)
	var nodeAsking uint64
	var length int
	var nRequested uint64

	// initialize InOutNeighbors
	for i := 1; i <= int(totalNodes); i++ {
		inNeighbors[uint64(i)] = []uint64{}
		outNeighbors[uint64(i)] = []uint64{}
	}

	fmt.Println("\n-------------------------------------------------------------")
	fmt.Println("------------------- SIMULATION ", time.Now().Format("01-02-2006 15:04:05"))
	fmt.Println("---------------------------------------------------------------")
	fmt.Println(" ")

	fmt.Println("totalNodes, Loop1, Loop2:", totalNodes, ",", loopNumber1, ",", loopNumber2)
	fmt.Println("zipfs, k, R, rho:", zipfs, ",", k, ",", bigR, ",", rho)

	/////////////////////////////////////////////////////////////////////
	/////////////////////////////////// FIRST ROUND
	/////////////////////////////////////////////////////////////////////

	fmt.Println("\n////////////////// FIRST ROUND")
	fmt.Println(" ")

	calculatePossibleNeigbors(listPossibleNeighbors, zipfsMana)
	// ToDo: alert if some nodes don't have possible neighbours

	// Create list of pairing prefereces by shuffling
	shuffleSlice(listPossibleNeighbors)
	//fmt.Println("ListPossibleNeighbors:", listPossibleNeighbors)

	for i := 1; i <= int(loopNumber1); i++ {
		nodeAsking = uint64(rand.Intn(int(totalNodes)) + 1)
		//outNeighborsAsking = outNeighbors[nodeAsking]
		length = len(listPossibleNeighbors[nodeAsking])
		nRequested = listPossibleNeighbors[nodeAsking][uint64(rand.Intn(length))]
		//inNeighboursCandidate = inNeighbors[candidate]

		updateInOutNeighbors(nodeAsking, nRequested, listPossibleNeighbors,
			inNeighbors, outNeighbors)
		// if !testUpdateInOut(inNeighbors, outNeighbors) {
		// 	break
		// }
	}
	// fmt.Println("ZipfsArray:", zipfsMana)
	//fmt.Println("InNeighbors:", inNeighbors)
	//fmt.Println("OutNeibours:", outNeighbors)

	nodesLowInDegree, nodesLowOutDegree, _, _ := calculateStatisticsLowDegreeNodes(inNeighbors, outNeighbors)

	////////////////////////////////////////////////////////////////////////////////
	///////////////////// SECOND ROUND: PAIRING LowOutDegreeN WITH LowInDegreeN
	////////////////////////////////////////////////////////////////////////////////

	fmt.Println(" ")
	if loopNumber2 > 0 {

		fmt.Println("////////////////// SECOND ROUND: PAIRING LowOutDegreeN WITH LowInDegreeN")
		fmt.Println(" ")

		var pLowInDN []uint64 // PossibleLowInDegreeNeighbors
		length = len(nodesLowOutDegree)

		for i := 0; i <= int(loopNumber2); i++ {
			nodeAsking = nodesLowOutDegree[rand.Intn(length)]
			pLowInDN = calculatePossibleLowInDegreeN(listPossibleNeighbors, inNeighbors, nodeAsking)
			if len(pLowInDN) > 0 {
				nRequested = pLowInDN[rand.Intn(len(pLowInDN))]
				updateInOutNeighborsLowPairing(nodeAsking, nRequested, listPossibleNeighbors,
					inNeighbors, outNeighbors, &nodesLowInDegree, &nodesLowOutDegree)

			}
		}

		//fmt.Println("InNeighbors:", inNeighbors)
		//fmt.Println("OutNeibours:", outNeighbors)

		calculateStatisticsLowDegreeNodes(inNeighbors, outNeighbors)
		fmt.Print(" ")
	}

}

///////////////////////////////////////////////////////////////
///////////////// DECLARATION OF FUNCTIONS
///////////////////////////////////////////////////////////////

func calculatePossibleLowInDegreeN(listPossibleNeighbors map[uint64][]uint64, inNeighbors map[uint64][]uint64, node uint64) (possibleLInDN []uint64) {
	possibleLInDN = make([]uint64, 0, totalNodes)
	for _, candidate := range listPossibleNeighbors[node] {
		if testLowInDegree(inNeighbors, candidate) {
			possibleLInDN = append(possibleLInDN, candidate)
		}
	}
	return
}

func testLowInDegree(inNeighbors map[uint64][]uint64, node uint64) bool {
	if len(inNeighbors[node]) >= int(inBound) {
		return false
	}
	return true
}

func testLowOutDegree(outNeighbors map[uint64][]uint64, node uint64) bool {
	if len(outNeighbors[node]) >= int(outBound) {
		return false
	}
	return true
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func testUpdateInOut(inNeighbors, outNeighbors map[uint64][]uint64) (test bool) {
	// In count
	var inCount int
	// for _, val := range inNeighbors {
	// 	if len(val) < int(inBound) {
	// 		inCount += int(inBound) - len(val)
	// 	}
	// }
	// // Out count
	// var outCount int
	// for _, val := range outNeighbors {
	// 	if len(val) < int(outBound) {
	// 		outCount += int(outBound) - len(val)
	// 	}
	// }
	for _, val := range inNeighbors {
		inCount += len(val)
	}
	// Out count
	var outCount int
	for _, val := range outNeighbors {
		outCount += len(val)
	}

	test = (inCount == outCount)
	if !test {
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println("inCount, outCount             :", inCount, outCount)
	}
	return
}

func calculateStatisticsLowDegreeNodes(inNeighbors, outNeighbors map[uint64][]uint64) (nodesLowInDegree, nodesLowOutDegree []uint64, distrInDegree map[uint64]uint64, distrOutDegree map[uint64]uint64) {

	nodesLowOutDegree = make([]uint64, 0, totalNodes)
	nodesLowInDegree = make([]uint64, 0, totalNodes)
	distrOutDegree = make(map[uint64]uint64, outBound+3)
	distrInDegree = make(map[uint64]uint64, inBound+3)

	// Calculate nodes with low In degree
	for i, val := range inNeighbors {
		if len(val) < int(inBound) {
			nodesLowInDegree = append(nodesLowInDegree, i)
			distrInDegree[uint64(len(val))]++
		} else {
			distrInDegree[uint64(inBound)]++
		}
	}

	// Calculate nodes with low Out degree
	for i, val := range outNeighbors {
		if len(val) < int(outBound) {
			nodesLowOutDegree = append(nodesLowOutDegree, i)
			distrOutDegree[uint64(len(val))]++
		} else {
			distrOutDegree[uint64(outBound)]++
		}
	}

	//////////////// TEST1 -- #edges by inDegree vs. #edges by outdegree
	var testIn, testOut int
	var test bool
	// LowIn count
	//fmt.Println("LowInDegreeNodes :", nodesLowInDegree)
	for _, node := range nodesLowInDegree {
		testIn += int(inBound) - len(inNeighbors[node])
	}
	// LowOut count
	//fmt.Println("LowOutDegreeNodes:", nodesLowOutDegree)
	for _, node := range nodesLowOutDegree {
		testOut += int(outBound) - len(outNeighbors[node])
	}
	// Test
	test = (testIn == testOut)

	////////////////////////////// TEST2
	var testIn1, testOut1 int
	var agree1, agree2, agree bool
	// LowInDistr count
	for degree, amount := range distrInDegree {
		testIn1 += (int(inBound) - int(degree)) * int(amount)
	}
	// LowOutDistr count
	for degree, amount := range distrOutDegree {
		testOut1 += (int(inBound) - int(degree)) * int(amount)
	}

	agree1 = (testIn1 == testIn)
	agree2 = (testOut1 == testOut)
	agree = (agree1 && agree2)

	// TEST 3
	// In count
	var inCount int
	for _, val := range inNeighbors {
		inCount += len(val)
	}
	// Out count
	var outCount int
	for _, val := range outNeighbors {
		outCount += len(val)
	}

	fmt.Println("#Nodes LowInDegree :", len(nodesLowInDegree))
	fmt.Println("#Nodes LowOutDegree:", len(nodesLowOutDegree))
	fmt.Println("DistrInDegree      :", distrInDegree)
	fmt.Println("DistrOutDegree     :", distrOutDegree)
	fmt.Println("InNeigh., OutNeigh.:", len(inNeighbors), len(outNeighbors))
	fmt.Println("Test1              :", test, "-- testIn:", testIn, "-- testOut:", testOut)
	fmt.Println("Test2              :", agree)
	fmt.Println("Test3              :", inCount == outCount)

	return
}

func updateInOutNeighborsLowPairing(nRequesting uint64, nAnswering uint64,
	mapPossibleNeighbors map[uint64][]uint64, inNeighbors map[uint64][]uint64,
	outNeighbors map[uint64][]uint64, nodesLowInDegree *[]uint64,
	nodesLowOutDegree *[]uint64) {

	preferencesRequester := mapPossibleNeighbors[nRequesting]
	preferencesNRequested := mapPossibleNeighbors[nAnswering]

	list := make([]uint64, 0, int(outBound)+1)
	preferences := make([]uint64, 0, totalNodes)

	lessThan := func(i, j int) bool {
		ranki, oki := find(preferences, list[i])
		rankj, okj := find(preferences, list[j])
		if oki && okj {
			return (ranki < rankj)
		} else if okj == true {
			return false
		} else {
			return true
		}
	}

	///////////// Updating OutNeighbors of Requester
	list = outNeighbors[nRequesting][:]
	preferences = preferencesRequester[:]
	list = append(list, nAnswering)
	list = removeDuplicateValues(list)
	numberOut := len(list)
	aboveBound1 := (numberOut > int(outBound))
	sort.SliceStable(list, lessThan)
	list1 := list[:]

	////////// Verify if Requester wants nRequested
	var agreementRequester bool = true
	if aboveBound1 && nAnswering == list1[int(outBound)] {
		agreementRequester = false
	}

	/////// Update InNeighbours of nRequested
	list = inNeighbors[nAnswering][:]
	preferences = preferencesNRequested[:]
	list = append(list, nRequesting)
	list = removeDuplicateValues(list)
	numberIn := len(list)
	aboveBound2 := (numberIn > int(inBound))
	sort.SliceStable(list, lessThan)
	list2 := list[:]

	//////////// Verify if nRequested wants Requester
	var agreementNAnswering bool = true
	if aboveBound2 && nRequesting == list2[int(inBound)] {
		agreementNAnswering = false
	}

	if agreementRequester && agreementNAnswering {

		//Update OutNeigh. of nRequesting
		cutPosition := min(numberOut, int(outBound))
		outNeighbors[nRequesting] = list1[:cutPosition]
		// Update nodesLowOutDegree
		updateLowDegreeList(outNeighbors, nodesLowOutDegree, nRequesting)

		//////////////////////////////// Update InNeighbors of DroppedNode
		if aboveBound1 {
			droppedNode := list1[int(outBound)]
			if droppedNode != nAnswering {
				list1 = inNeighbors[droppedNode][:]
				pos, ok := find(list1, nRequesting)
				//fmt.Println("postition: ", pos)
				if ok {
					inNeighbors[droppedNode] = removeIndex(list1, pos)
					//inNeighbors[droppedNode] = append(list1[:int(pos)], list1[int(pos)+1:]...)
					// Update NLowInDegree
					updateLowDegreeList(inNeighbors, nodesLowInDegree, droppedNode)
				}
			}

		}

		// Update InNeigh. of nAnswering
		cutPosition = min(numberIn, int(inBound))
		inNeighbors[nAnswering] = list2[:cutPosition]
		// Update NLowInDegree
		updateLowDegreeList(inNeighbors, nodesLowInDegree, nAnswering)

		//////////////////// Update InNeibours of Dropped node
		if aboveBound2 {
			droppedNode := list2[int(inBound)]
			if droppedNode != nRequesting {
				list2 = outNeighbors[droppedNode][:]
				pos, ok := find(list2, nAnswering)
				if ok {
					outNeighbors[droppedNode] = removeIndex(list2, pos)
					//outNeighbors[droppedNode] = append(list2[:int(pos)], list2[int(pos)+1:]...)
					// Update NLowOutDegree
					updateLowDegreeList(outNeighbors, nodesLowOutDegree, droppedNode)
				}
			}
		}

	}

	return

}

func updateLowDegreeList(inOutNeighbors map[uint64][]uint64, listLowInOutDegree *[]uint64, node uint64) {
	pos, ok := find(*listLowInOutDegree, node)
	test := testLowInDegree(inOutNeighbors, node)
	if (ok == false) && test {
		*listLowInOutDegree = append(*listLowInOutDegree, node)
		return
	} else if ok && (test == false) {
		removeIndex(*listLowInOutDegree, pos)
		return
	}
	return
}

func removeIndex(s []uint64, index uint64) []uint64 {
	return append(s[:index], s[index+1:]...)
}

func updateInOutNeighbors(nRequesting uint64, nAnswering uint64,
	mapPossibleNeighbors map[uint64][]uint64, inNeighbors map[uint64][]uint64,
	outNeighbors map[uint64][]uint64) {

	preferencesRequester := mapPossibleNeighbors[nRequesting]
	preferencesNRequested := mapPossibleNeighbors[nAnswering]

	list := make([]uint64, 0, int(outBound)+1)
	preferences := make([]uint64, 0, totalNodes)

	lessThan := func(i, j int) bool {
		ranki, oki := find(preferences, list[i])
		rankj, okj := find(preferences, list[j])
		if oki && okj {
			return (ranki < rankj)
		} else if okj == true {
			return false
		} else {
			return true
		}
	}

	///////////// Updating OutNeighbors of Requester
	list = outNeighbors[nRequesting][:]
	preferences = preferencesRequester[:]
	list = append(list, nAnswering)
	list = removeDuplicateValues(list)
	numberOut := len(list)
	aboveBound1 := (numberOut > int(outBound))
	sort.SliceStable(list, lessThan)
	list1 := list[:]

	////////// Verify if Requester wants nRequested
	var agreementRequester bool = true
	if aboveBound1 && nAnswering == list1[int(outBound)] {
		agreementRequester = false
	}

	/////// Update InNeighbours of nRequested
	list = inNeighbors[nAnswering][:]
	preferences = preferencesNRequested[:]
	list = append(list, nRequesting)
	list = removeDuplicateValues(list)
	numberIn := len(list)
	aboveBound2 := (numberIn > int(inBound))
	sort.SliceStable(list, lessThan)
	list2 := list[:]

	//////////// Verify if nRequested wants Requester
	var agreementNAnswering bool = true
	if aboveBound2 && nRequesting == list2[int(inBound)] {
		agreementNAnswering = false
	}

	if agreementRequester && agreementNAnswering {
		var droppedNode1, droppedNode2 uint64
		var list3, list4 []uint64

		//Update OutNeigh. of nRequesting
		cutPosition1 := min(numberOut, int(outBound))
		outNeighbors[nRequesting] = list1[:cutPosition1]
		//////////////////////////////// Update InNeighbors of DroppedNode
		if aboveBound1 {
			droppedNode1 = list1[int(outBound)]
			//if droppedNode != nAnswering {
			list3 = inNeighbors[droppedNode1]
			pos, ok := find(list3, nRequesting)
			//fmt.Println("postition: ", pos)
			if ok {
				inNeighbors[droppedNode1] = removeIndex(list3, pos)
				//inNeighbors[droppedNode] = append(list1[:int(pos)], list1[int(pos)+1:]...)
			}
			//}

		}

		// Update InNeigh. of nAnswering
		cutPosition2 := min(numberIn, int(inBound))
		inNeighbors[nAnswering] = list2[:cutPosition2]
		//////////////////// Update InNeibours of Dropped node
		if aboveBound2 {
			droppedNode2 = list2[int(inBound)]
			//if droppedNode != nRequesting {
			list4 = outNeighbors[droppedNode2]
			pos, ok := find(list4, nAnswering)
			if ok {
				outNeighbors[droppedNode2] = removeIndex(list4, pos)
				//outNeighbors[droppedNode] = append(list2[:int(pos)], list2[int(pos)+1:]...)
			}
			//}
		}

		// Debugging
		if false && !testUpdateInOut(inNeighbors, outNeighbors) {
			fmt.Println("nAnswering, nAsking       :", nAnswering, nRequesting)
			fmt.Println("InNeighNAnsw, OutNeighNAsk:", inNeighbors[nAnswering], outNeighbors[nRequesting])
			fmt.Println("AgreementReached          :", agreementRequester && agreementNAnswering)
			fmt.Println("DroppedNode1, DroppedNode2:", droppedNode1, droppedNode2)
			fmt.Println("AboveBound1, AboveBound2  :", aboveBound1, aboveBound2)
			fmt.Println("list1                     :", list1)
			fmt.Println("list2                     :", list2)
			fmt.Println("list3                     :", list3)
			fmt.Println("list4                     :", list4)
			fmt.Println("")
		}

	}

	return
}

func find(source []uint64, value uint64) (uint64, bool) {
	for i, item := range source {
		if item == value {
			return uint64(i), true
		}
	}
	return 0, false
}

func calculatePossibleNeigbors(listPossibleNeighbors map[uint64][]uint64,
	mana map[uint64]uint64) {
	for i := 1; i <= int(totalNodes); i++ {
		for j := i + 1; j <= int(totalNodes); j++ {
			// Condition1 := mana[i] < mana[j]*rho
			if mana[uint64(i)] < mana[uint64(j)]*rho {
				//fmt.Println("i,j:", i, j)
				listPossibleNeighbors[uint64(i)] = append(listPossibleNeighbors[uint64(i)], uint64(j))
				listPossibleNeighbors[uint64(j)] = append(listPossibleNeighbors[uint64(j)], uint64(i))
			} else if uint64(math.Abs(float64(i-j))) <= bigR {
				//fmt.Println("i,j:", i, j)
				listPossibleNeighbors[uint64(i)] = append(listPossibleNeighbors[uint64(i)], uint64(j))
				listPossibleNeighbors[uint64(j)] = append(listPossibleNeighbors[uint64(j)], uint64(i))
			} else {
				break
			}
		}
	}
	return
}

func shuffleSlice(in map[uint64][]uint64) {
	//rand.Seed(time.Now().Unix())
	for _, list := range in {
		rand.Shuffle(len(list), func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})
	}
	return
}

func removeDuplicateValues(intSlice []uint64) []uint64 {
	keys := make(map[uint64]bool)
	list := []uint64{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func createZipfsMana(zipfsMana map[uint64]uint64) {
	scalingFactor := math.Pow(10, 10)
	for i := 1; i < int(totalNodes+1); i++ {
		zipfsMana[uint64(i)] = uint64(math.Pow(float64(i), -zipfs) * scalingFactor)
	}
	return
}
