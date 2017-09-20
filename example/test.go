package main

import (
	"../../btlib"
	"fmt"
)




func main() {
	pair := "BTC-NEO"
	bt := btlib.Btlib{}.NewClient()
	candles, err := bt.GetCandles(pair, "hour")
		if err != nil {
			fmt.Printf("%s - Get Candle\n", pair)
		} else {
			haCancles, _ := bt.HeikinAshi(candles)
			bb := bt.BB(haCancles, 30, 2)
			for _, b := range bb {
				fmt.Printf("%#v\n", b)
			}
		}
}



//var pairs = []string{"BTC-1ST", "BTC-2GIVE", "BTC-ABY", "BTC-ADT", "BTC-ADX", "BTC-AEON", "BTC-AGRS", "BTC-AMP", "BTC-ANT", "BTC-APX", "BTC-ARDR", "BTC-ARK", "BTC-AUR", "BTC-BAT", "BTC-BAY", "BTC-BCC", "BTC-BCY", "BTC-BITB", "BTC-BLITZ", "BTC-BLK", "BTC-BLOCK", "BTC-BNT", "BTC-BRK", "BTC-BRX", "BTC-BSD", "BTC-BTA", "BTC-BTCD", "BTC-BTS", "BTC-BURST", "BTC-BYC", "BTC-CANN", "BTC-CFI", "BTC-CLAM", "BTC-CLOAK", "BTC-CLUB", "BTC-COVAL", "BTC-CPC", "BTC-CRB", "BTC-CRW", "BTC-CURE", "BTC-CVC", "BTC-DAR", "BTC-DASH", "BTC-DCR", "BTC-DCT", "BTC-DGB", "BTC-DGD", "BTC-DMD", "BTC-DOGE", "BTC-DOPE", "BTC-DRACO", "BTC-DTB", "BTC-DYN", "BTC-EBST", "BTC-EDG", "BTC-EFL", "BTC-EGC", "BTC-EMC", "BTC-EMC2", "BTC-ENRG", "BTC-ERC", "BTC-ETC", "BTC-ETH", "BTC-EXCL", "BTC-EXP", "BTC-FAIR", "BTC-FCT", "BTC-FLDC", "BTC-FLO", "BTC-FTC", "BTC-FUN", "BTC-GAM", "BTC-GAME", "BTC-GBG", "BTC-GBYTE", "BTC-GCR", "BTC-GEO", "BTC-GLD", "BTC-GNO", "BTC-GNT", "BTC-GOLOS", "BTC-GRC", "BTC-GRS", "BTC-GUP", "BTC-HMQ", "BTC-INCNT", "BTC-INFX", "BTC-IOC", "BTC-ION", "BTC-IOP", "BTC-KMD", "BTC-KORE", "BTC-LBC", "BTC-LGD", "BTC-LMC", "BTC-LSK", "BTC-LTC", "BTC-LUN", "BTC-MAID", "BTC-MCO", "BTC-MEME", "BTC-MLN", "BTC-MONA", "BTC-MTL", "BTC-MUE", "BTC-MUSIC", "BTC-MYST", "BTC-NAV", "BTC-NBT", "BTC-NEO", "BTC-NEOS", "BTC-NLG", "BTC-NMR", "BTC-NXC", "BTC-NXS", "BTC-NXT", "BTC-OK", "BTC-OMG", "BTC-OMNI", "BTC-PART", "BTC-PAY", "BTC-PDC", "BTC-PINK", "BTC-PIVX", "BTC-PKB", "BTC-POT", "BTC-PPC", "BTC-PTC", "BTC-PTOY", "BTC-QRL", "BTC-QTUM", "BTC-QWARK", "BTC-RADS", "BTC-RBY", "BTC-RDD", "BTC-REP", "BTC-RISE", "BTC-RLC", "BTC-SAFEX", "BTC-SBD", "BTC-SC", "BTC-SEQ", "BTC-SHIFT", "BTC-SIB", "BTC-SLR", "BTC-SLS", "BTC-SNGLS", "BTC-SNRG", "BTC-SNT", "BTC-SPHR", "BTC-SPR", "BTC-START", "BTC-STEEM", "BTC-STORJ", "BTC-STRAT", "BTC-SWIFT", "BTC-SWT", "BTC-SYNX", "BTC-SYS", "BTC-THC", "BTC-TIME", "BTC-TKN", "BTC-TKS", "BTC-TRIG", "BTC-TRST", "BTC-TRUST", "BTC-TX", "BTC-UBQ", "BTC-UNB", "BTC-UNO", "BTC-VIA", "BTC-VOX", "BTC-VRC", "BTC-VRM", "BTC-VTC", "BTC-VTR", "BTC-WAVES", "BTC-WINGS", "BTC-XAUR", "BTC-XCP", "BTC-XDN", "BTC-XEL", "BTC-XEM", "BTC-XLM", "BTC-XMG", "BTC-XMR", "BTC-XMY", "BTC-XRP", "BTC-XST", "BTC-XVC", "BTC-XVG", "BTC-XWC", "BTC-XZC", "BTC-ZCL", "BTC-ZEC", "BTC-ZEN"}
//func main()  {
//	bt := btlib.Btlib{}.NewClient()
//	for _, pair := range pairs {
//		candles, err := bt.GetCandles(pair, "hour")
//
//		if err != nil {
//			fmt.Printf("%s - Get Candle\n", pair)
//		} else {
//			haCancles, _ := bt.HeikinAshi(candles)
//			vols, _ := bt.Vol(haCancles, 20)
//			if haCancles[len(haCancles)-1].BaseVolume > vols[len(vols)-1] {
//				if haCancles[len(haCancles)-1].Close > haCancles[len(haCancles)-1].Open {
//					fmt.Printf("%s - Buy Vol\n", pair)
//				} else {
//					//fmt.Printf("%s - Sell Vol\n", pair)
//				}
//
//			}
//		}
//
//	}
//}