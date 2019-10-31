/*
Copyright © 2019 Steven Truong struong996@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// variables for flags
var cfgFile string
var n int8
var s []float32

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bill",
	Short: "A CLI tool that calculates how much you need to pay for your bill.",
	Long:  ``,
	Args:  validateBillArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: BillImpl,
}

// BillImpl contains the logic of this CLI tool.
func BillImpl(cmd *cobra.Command, args []string) {
	// 1. figure out the principle
	var principle float32
	if len(args) > 0 {
		parsedFloat, err := strconv.ParseFloat(args[0], 32)
		if err != nil {
			log.Fatal(err)
		}
		principle = float32(parsedFloat)
	}
	// 2. calculate the amount to pay from the split
	// check if a split is specified
	var split float32 = 0
	if len(s) > 0 {
		for i := range s {
			split += s[i]
		}
		split = split / float32(n)
	}
	// 3. add the result from (1) and (2)
	result := principle + split

	fmt.Println("The amount to pay is: ", result)
}

// validateBillArgs is the argument validator for this CLI tool.
func validateBillArgs(cmd *cobra.Command, args []string) error {
	// if no principle is specified, s should be specified
	if len(args) == 0 && len(s) == 0 {
		return errors.New("Please provide a principle amount or list of items to split for the bill")
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bill.yaml)")
	// (bevensteven) added flags for tool
	rootCmd.Flags().Int8VarP(&n, "numPeople", "n", 1, "The number of people to split the to divide arguments with.")
	rootCmd.Flags().Float32SliceVarP(&s, "splits", "s", make([]float32, 0), `An iterable of values to be split among n people. Example: -s 10,10,10`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bill" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bill")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
