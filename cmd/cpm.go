package cmd

// var cpmCmd = &cobra.Command{
// 	Use:   "cpm <rom>",
// 	Args:  cobra.ExactArgs(1),
// 	Short: "Run CP/M",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		logger := newDefaultLogger()
// 		if debug {
// 			logger.SetLevel(log.DebugLevel)
// 		}

// 		fileName := filepath.Base(args[0])
// 		data, err := os.ReadFile(args[0])
// 		if err != nil {
// 			logger.Fatal(err)
// 		}

// 		// Assuming NewCPMHardware() sets up the CP/M environment
// 		vm := emulator.NewCPU8080(&data, startAddress, NewCPMHardware())
// 		vm.StartInterruptRoutines()
// 		vm.Logger = logger
// 		vm.Options.UnlimitedTPS = true

// 		ebiten.SetWindowTitle(fileName)
// 		if vm.Options.UnlimitedTPS {
// 			ebiten.SetTPS(ebiten.SyncWithFPS)
// 		} else {
// 			ebiten.SetTPS(60)
// 		}
// 		ebiten.SetWindowSize(vm.Hardware.Width()*vm.Hardware.Scale(), vm.Hardware.Height()*vm.Hardware.Scale())

// 		if err := ebiten.RunGame(vm); err != nil && err != ebiten.Termination {
// 			logger.Fatal(err)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(cpmCmd)
// }
