package cryptonote

func DecomposeAmountIntoDigits(amount, dustThreshold uint64) ([]uint64, []uint64) {
	var dusts []uint64
	var chunks []uint64

	if amount != 0 {
		isDustHandled := false
		order := uint64(1)
		dust := uint64(0)

		for amount != 0 {
			chunk := (amount % 10) * order
			amount /= 10
			order *= 10

			if dust+chunk <= dustThreshold {
				dust += chunk
			} else {
				if !isDustHandled && dust != 0 {
					dusts = append(dusts, dust)
					isDustHandled = true
				}
				if chunk != 0 {
					chunks = append(chunks, chunk)
				}
			}
		}

		if !isDustHandled && dust != 0 {
			dusts = append(dusts, dust)
		}
	}

	return chunks, dusts
}
