package main

import (
	"testing"
	"time"
)

func TestTimeSeriesIpDuplicates_Add_DuplicateWithinSlot(t *testing.T) {
	duplicatesChan := make(chan string, 1)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.1", time.Now().Add(5*time.Second))

	select {
	case duplicateIP := <-duplicatesChan:
		if duplicateIP != "192.168.0.1" {
			t.Errorf("Expected duplicate IP: %s, but got: %s", "192.168.0.1", duplicateIP)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout: Expected duplicate IP to be sent to the channel, but no IP was received.")
	}
}

func TestTimeSeriesIpDuplicates_Add_DuplicateAfterSlot(t *testing.T) {
	duplicatesChan := make(chan string, 1)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.1", time.Now().Add(15*time.Second))
	ts.Add("192.168.0.1", time.Now().Add(25*time.Second))

	select {
	case duplicateIP := <-duplicatesChan:
		t.Errorf("Received duplicate IP: %s, but didn't expect any duplicates. IP: %s", duplicateIP, "192.168.0.1")
	case <-time.After(1 * time.Second):
		// No duplicate IP should be received.
	}
}

func TestTimeSeriesIpDuplicates_Add_ExceedThreshold(t *testing.T) {
	duplicatesChan := make(chan string, 1)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.2", time.Now())
	ts.Add("192.168.0.2", time.Now().Add(2*time.Second))
	ts.Add("192.168.0.3", time.Now().Add(8*time.Second))

	select {
	case duplicateIP := <-duplicatesChan:
		if duplicateIP != "192.168.0.2" {
			t.Errorf("Expected duplicate IP: %s, but got: %s", "192.168.0.2", duplicateIP)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout: Expected duplicate IP to be sent to the channel, but no IP was received.")
	}
}

func TestTimeSeriesIpDuplicates_Add_NoDuplicate(t *testing.T) {
	duplicatesChan := make(chan string)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.3", time.Now())

	select {
	case duplicateIP := <-duplicatesChan:
		t.Errorf("Received duplicate IP: %s, but didn't expect any duplicates. IP: %s", duplicateIP, "192.168.0.3")
	case <-time.After(1 * time.Second):
		// No duplicate IP should be received.
	}
}

func TestTimeSeriesIpDuplicates_Add_MultipleSlots(t *testing.T) {
	duplicatesChan := make(chan string, 2)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	// Add duplicates within the first time slot.
	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.1", time.Now().Add(5*time.Second))

	// Sleep for 11 seconds to create a new time slot.
	time.Sleep(11 * time.Second)

	// Add duplicates within the second time slot.
	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.1", time.Now().Add(5*time.Second))

	select {
	case duplicateIP := <-duplicatesChan:
		if duplicateIP != "192.168.0.1" {
			t.Errorf("Expected duplicate IP: %s, but got: %s", "192.168.0.1", duplicateIP)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout: Expected duplicate IP to be sent to the channel, but no IP was received.")
	}
}

func TestTimeSeriesIpDuplicates_Add_DifferentIPs(t *testing.T) {
	duplicatesChan := make(chan string)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.2", time.Now())

	// No duplicates for either IP, so we shouldn't receive anything on the channel.
	select {
	case duplicateIP := <-duplicatesChan:
		t.Errorf("Received duplicate IP: %s, but didn't expect any duplicates. IP: %s", duplicateIP, "192.168.0.1")
	case <-time.After(1 * time.Second):
		// No duplicate IP should be received.
	}
}

func TestTimeSeriesIpDuplicates_Add_DifferentThreshold(t *testing.T) {
	duplicatesChan := make(chan string)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 3)
	defer close(duplicatesChan)

	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.1", time.Now().Add(5*time.Second))

	// The threshold is 3, so we shouldn't receive anything on the channel.
	select {
	case duplicateIP := <-duplicatesChan:
		t.Errorf("Received duplicate IP: %s, but didn't expect any duplicates. IP: %s", duplicateIP, "192.168.0.1")
	case <-time.After(1 * time.Second):
		// No duplicate IP should be received.
	}
}

func TestTimeSeriesIpDuplicates_Add_DuplicateAfterThreshold(t *testing.T) {
	duplicatesChan := make(chan string, 1)
	ts := NewTimeSeriesIpDuplicates(time.Now(), duplicatesChan, 10*time.Second, 2)
	defer close(duplicatesChan)

	ts.Add("192.168.0.1", time.Now())
	ts.Add("192.168.0.1", time.Now().Add(11*time.Second))

	// The third occurrence comes after the time slot, so it shouldn't be treated as a duplicate.
	select {
	case duplicateIP := <-duplicatesChan:
		t.Errorf("Received duplicate IP: %s, but didn't expect any duplicates. IP: %s", duplicateIP, "192.168.0.1")
	case <-time.After(1 * time.Second):
		// No duplicate IP should be received.
	}
}
