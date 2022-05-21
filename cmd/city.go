/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zer0-star/naro-go-database/city"
)

// cityCmd represents the city command
var cityCmd = &cobra.Command{
	Use:   "city",
	Short: "Print the population of city",
	Long:  `Print the population of city and percentage of city population to country population`,
	Run:   city.Run,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(cityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
