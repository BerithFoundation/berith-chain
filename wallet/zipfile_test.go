package main

import (
	"log"
	"testing"
)

func TestZipEncrypted(t *testing.T) {
	directory := "C:\\Users\\Usman\\zipTest"
	zippedFile := "C:\\Users\\Usman\\zipTest\\Zipped.zip"
	if err := ZipSecure(directory, zippedFile, "usman"); err != nil {
		log.Fatalln(err)
	}
}

func TestUnzipSecureRightPassword(t *testing.T) {
	zippedFile := "C:\\Users\\Usman\\zipTest\\Zipped.zip"
	outputDirectory := "C:\\Users\\Usman\\zipTest\\Zipped"
	if err := UnzipSecure(zippedFile, outputDirectory, "usman"); err != nil {
		log.Fatalln(err)
	}
}

func TestUnzipSecureWrongPassword(t *testing.T) {
	zippedFile := "C:\\Users\\Usman\\zipTest\\Zipped.zip"
	outputDirectory := "C:\\Users\\Usman\\zipTest\\Zipped"
	if err := UnzipSecure(zippedFile, outputDirectory, "abc"); err != nil {
		log.Fatalln(err)
	}
}
