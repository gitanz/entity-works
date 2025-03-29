package main

import "testing"

func TestSaysHello(t *testing.T) {
	expectedGreeting := "Hello"
	greeting := sayHello();
	if greeting !=  expectedGreeting  {
		t.Errorf("Expected %v, got  %v", expectedGreeting, greeting)
	}
}