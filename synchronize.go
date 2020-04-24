package main

import (
	"fmt"
	"io"

	"github.com/golang-collections/collections/set"
)


func synchronize(store1 ObjectStore, store2 ObjectStore) {

	fmt.Println("obtaining infosTarget...")
	infosTargetCh := store2.GetInfos()

	fmt.Println("obtaining infosSource...")
	infosSourceCh := store1.GetInfos()

	setTarget := set.New()
	for m := range infosTargetCh {
		setTarget.Insert(m)
	}
	fmt.Printf("obtained %v infosTarget.\n", setTarget.Len())

	setSource := set.New()
	for m := range infosSourceCh {
		setSource.Insert(m)
	}
	fmt.Printf("obtained %v infosSource.\n", setSource.Len())

	diffAdd := setSource.Difference(setTarget)
	fmt.Printf("obtained %v additive differences \n", diffAdd.Len())

	diffSub := setTarget.Difference(setSource)
	fmt.Printf("obtained %v subtractive differences \n", diffSub.Len())
	
	var diffAddInfos []ObjectInfo
	diffAdd.Do(func(info interface{}) {
		diffAddInfos = append(diffAddInfos, info.(ObjectInfo))
	})

	lenDiff := len(diffAddInfos)
	for i, diffInfo := range diffAddInfos {
		address := diffInfo.Address

		writer, err := store2.GetWriter(address)
		if err != nil {
			fmt.Printf("Error: %v \n", err)
			continue
		}
		defer writer.Close()

		reader, err := store1.GetReader(address)
		if err != nil {
			fmt.Printf("Error: %v \n", err)
			continue
		}
		defer reader.Close()

		fmt.Printf("%v of %v | Copy %+v \n", i+1, lenDiff, address)
		_, err = io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("Error: %v \n", err)
		}
	}
}
