package pointers

import "testing"

func TestWallet(t *testing.T) {
	t.Run("Wallet should update balance for deposit", func(t *testing.T) {
		w := Wallet{}
		w.Deposit(Bitcoin(10.0))
		got := w.Balance()
		expected := Bitcoin(10.0)

		assertFloatEqual(t, got, expected)
	})

	t.Run("Wallet should update balance for withdraw", func(t *testing.T) {
		w := Wallet{10.0}
		w.Withdraw(Bitcoin(5.0))
		got := w.Balance()
		expected := Bitcoin(5.0)

		assertFloatEqual(t, got, expected)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		w := Wallet{10.0}
		err := w.Withdraw(Bitcoin(15.0))

		got := w.Balance()
		expected := Bitcoin(10.0)

		assertError(t, err, ErrWithdraw)
		assertFloatEqual(t, got, expected)
	})
}

func assertFloatEqual(t *testing.T, got Bitcoin, expected Bitcoin) {
	t.Helper()

	if got != expected {
		t.Errorf("Wrong value: got %s, expected %s", got, expected)
	}
}

func assertError(t *testing.T, err error, expectedErr error) {
	if err == nil {
		t.Fatal("Wanted an error but didn't get one")
	}

	if err != expectedErr {
		t.Errorf("wrong error message")
	}
}
