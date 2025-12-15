package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	baseURL          = "http://localhost:8080"
	chromeDriverPort = 4444
)

var wd selenium.WebDriver

func setupWebDriver(t *testing.T) {
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	chromeCaps := chrome.Capabilities{
		Args: []string{
			// "--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--disable-cache",
			"--disable-application-cache",
			"--disable-offline-load-stale-cache",
			"--window-size=1920,1080",
		},
		Path:            "",
		ExcludeSwitches: []string{},
		Extensions:      []string{},
		LocalState:      map[string]interface{}{},
		Prefs:           map[string]interface{}{},
		Detach:          new(bool),
		DebuggerAddr:    "",
		MinidumpPath:    "",
		MobileEmulation: &chrome.MobileEmulation{},
		WindowTypes:     []string{},
		AndroidPackage:  "",
		W3C:             false,
	}
	caps.AddChrome(chromeCaps)

	var err error
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", chromeDriverPort))
	if err != nil {
		t.Fatalf("Failed to create WebDriver: %v", err)
	}
}

func teardownWebDriver() {
	if wd != nil {
		wd.Quit()
	}
}

func waitForElement(t *testing.T, by, value string, timeout time.Duration) selenium.WebElement {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		elem, err := wd.FindElement(by, value)
		if err == nil {
			return elem
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Element not found: %s=%s", by, value)
	return nil
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func waitForNotificationWithText(t *testing.T, expectedText string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	var lastText string
	for time.Now().Before(deadline) {
		notif, err := wd.FindElement(selenium.ByID, "notification")
		if err == nil {
			text, _ := notif.Text()
			className, _ := notif.GetAttribute("class")
			lastText = text
			if text == expectedText && !contains(className, "hidden") {
				return
			}
			if !contains(className, "hidden") && text != "" {
				time.Sleep(100 * time.Millisecond)
				text, _ = notif.Text()
				if text == expectedText {
					return
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Expected notification '%s', but got '%s'", expectedText, lastText)
}

func waitForNotificationToDisappear(t *testing.T, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		notif, err := wd.FindElement(selenium.ByID, "notification")
		if err == nil {
			className, _ := notif.GetAttribute("class")
			if contains(className, "hidden") {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func loginHelper(t *testing.T) {
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())
	password := "testpass123"

	// Register first - force refresh to clear cache
	if err := wd.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}
	wd.Refresh() // Hard refresh to reload JS files

	registerTab := waitForElement(t, selenium.ByXPATH, "//button[@data-auth-tab='register']", 5*time.Second)
	registerTab.Click()
	time.Sleep(500 * time.Millisecond)

	usernameInput := waitForElement(t, selenium.ByID, "register-username", 2*time.Second)
	usernameInput.SendKeys(username)

	emailInput := waitForElement(t, selenium.ByID, "register-email", 1*time.Second)
	emailInput.SendKeys(username + "@test.com")

	passwordInput := waitForElement(t, selenium.ByID, "register-password", 1*time.Second)
	passwordInput.SendKeys(password)

	passwordConfirmInput := waitForElement(t, selenium.ByID, "register-password-confirm", 1*time.Second)
	passwordConfirmInput.SendKeys(password)

	acceptTerms, err := wd.FindElement(selenium.ByID, "accept-terms")
	if err != nil {
		t.Fatalf("Failed to find accept-terms checkbox: %v", err)
	}
	acceptTerms.Click()

	submitBtn, err := wd.FindElement(selenium.ByID, "register-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find register submit button: %v", err)
	}
	submitBtn.Click()

	time.Sleep(3 * time.Second)

	appContainer := waitForElement(t, selenium.ByID, "app-container", 10*time.Second)
	className, _ := appContainer.GetAttribute("class")
	if contains(className, "hidden") {
		authContainer, _ := wd.FindElement(selenium.ByID, "auth-container")
		authClass, _ := authContainer.GetAttribute("class")
		t.Fatalf("Login helper failed: app container still hidden. Auth container: %s", authClass)
	}
	time.Sleep(1 * time.Second)
}

// Test 1: User Registration
func TestUserRegistration(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	username := fmt.Sprintf("user_%d", time.Now().Unix())

	if err := wd.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	registerTab := waitForElement(t, selenium.ByXPATH, "//button[@data-auth-tab='register']", 5*time.Second)
	registerTab.Click()
	time.Sleep(500 * time.Millisecond)

	usernameInput := waitForElement(t, selenium.ByID, "register-username", 2*time.Second)
	usernameInput.SendKeys(username)

	emailInput := waitForElement(t, selenium.ByID, "register-email", 1*time.Second)
	emailInput.SendKeys(username + "@test.com")

	passwordInput := waitForElement(t, selenium.ByID, "register-password", 1*time.Second)
	passwordInput.SendKeys("password123")

	passwordConfirmInput := waitForElement(t, selenium.ByID, "register-password-confirm", 1*time.Second)
	passwordConfirmInput.SendKeys("password123")

	acceptTerms, err := wd.FindElement(selenium.ByID, "accept-terms")
	if err != nil {
		t.Fatalf("Failed to find accept-terms checkbox: %v", err)
	}
	acceptTerms.Click()

	submitBtn, err := wd.FindElement(selenium.ByID, "register-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find register submit button: %v", err)
	}
	submitBtn.Click()

	time.Sleep(3 * time.Second)

	appContainer := waitForElement(t, selenium.ByID, "app-container", 10*time.Second)
	className, _ := appContainer.GetAttribute("class")
	if contains(className, "hidden") {
		authContainer, _ := wd.FindElement(selenium.ByID, "auth-container")
		authClass, _ := authContainer.GetAttribute("class")
		t.Fatalf("App container hidden after registration. Auth container: %s", authClass)
	}
}

// Test 2: User Login
func TestUserLogin(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	username := fmt.Sprintf("user_%d", time.Now().Unix())

	if err := wd.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	// Register first
	registerTab := waitForElement(t, selenium.ByXPATH, "//button[@data-auth-tab='register']", 5*time.Second)
	registerTab.Click()
	time.Sleep(500 * time.Millisecond)

	usernameInput := waitForElement(t, selenium.ByID, "register-username", 2*time.Second)
	usernameInput.SendKeys(username)

	emailInput := waitForElement(t, selenium.ByID, "register-email", 1*time.Second)
	emailInput.SendKeys(username + "@test.com")

	passwordInput := waitForElement(t, selenium.ByID, "register-password", 1*time.Second)
	passwordInput.SendKeys("password123")

	passwordConfirmInput := waitForElement(t, selenium.ByID, "register-password-confirm", 1*time.Second)
	passwordConfirmInput.SendKeys("password123")

	acceptTerms, err := wd.FindElement(selenium.ByID, "accept-terms")
	if err != nil {
		t.Fatalf("Failed to find accept-terms checkbox: %v", err)
	}
	acceptTerms.Click()

	submitBtn, err := wd.FindElement(selenium.ByID, "register-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find register submit button: %v", err)
	}
	submitBtn.Click()

	waitForElement(t, selenium.ByID, "app-container", 5*time.Second)

	// Logout
	userMenuBtn := waitForElement(t, selenium.ByID, "user-menu-btn", 2*time.Second)
	userMenuBtn.Click()
	time.Sleep(500 * time.Millisecond)

	logoutBtn := waitForElement(t, selenium.ByID, "logout-menu-item", 2*time.Second)
	logoutBtn.Click()
	time.Sleep(500 * time.Millisecond)

	// Login
	loginTab := waitForElement(t, selenium.ByXPATH, "//button[@data-auth-tab='login']", 2*time.Second)
	loginTab.Click()
	time.Sleep(500 * time.Millisecond)

	loginUsername := waitForElement(t, selenium.ByID, "login-username", 2*time.Second)
	loginUsername.SendKeys(username)

	loginPassword := waitForElement(t, selenium.ByID, "login-password", 1*time.Second)
	loginPassword.SendKeys("password123")

	loginSubmit, err := wd.FindElement(selenium.ByID, "login-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find login submit button: %v", err)
	}
	loginSubmit.Click()

	time.Sleep(3 * time.Second)

	appContainer := waitForElement(t, selenium.ByID, "app-container", 10*time.Second)
	className, _ := appContainer.GetAttribute("class")
	if contains(className, "hidden") {
		authContainer, _ := wd.FindElement(selenium.ByID, "auth-container")
		authClass, _ := authContainer.GetAttribute("class")
		t.Fatalf("App container hidden after login. Auth container: %s", authClass)
	}
}

// Test 3: Create Book
func TestCreateBook(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
	addBookBtn.Click()

	// Step 1: Fill title
	titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.SendKeys("Test Book Title")

	// Click Next to go to step 2
	nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	// Step 2: Fill description
	descInput := waitForElement(t, selenium.ByID, "book-description", 1*time.Second)
	descInput.SendKeys("Test book description")

	// Click Next to go to step 3 (review)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	// Step 3: Submit
	submitBtn, err := wd.FindElement(selenium.ByID, "book-form-submit")
	if err != nil {
		t.Fatalf("Failed to find book submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Book created successfully", 5*time.Second)
}

// Test 4: Update and Delete Book
func TestUpdateDeleteBook(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	// Create a book first
	addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
	addBookBtn.Click()

	// Step 1: Fill title
	titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.SendKeys("Book to Update")

	// Navigate through steps to submit
	nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	submitBtn, err := wd.FindElement(selenium.ByID, "book-form-submit")
	if err != nil {
		t.Fatalf("Failed to find book submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Book created successfully", 5*time.Second)
	waitForNotificationToDisappear(t, 10*time.Second)

	// Update the book
	editBtns, err := wd.FindElements(selenium.ByClassName, "btn-primary")
	if err != nil || len(editBtns) == 0 {
		t.Fatalf("Failed to find edit buttons: %v", err)
	}
	for _, btn := range editBtns {
		text, _ := btn.Text()
		if text == "Edit" {
			btn.Click()
			time.Sleep(500 * time.Millisecond)
			break
		}
	}

	titleInput = waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.Clear()
	titleInput.SendKeys("Updated Book Title")

	// Navigate through steps to submit
	nextBtn = waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	submitBtn, err = wd.FindElement(selenium.ByID, "book-form-submit")
	if err != nil {
		t.Fatalf("Failed to find book submit button for update: %v", err)
	}
	submitBtn.Click()

	// Wait for async book reload to complete
	time.Sleep(3 * time.Second)

	waitForNotificationWithText(t, "Book updated successfully", 10*time.Second)

	// Delete the book
	time.Sleep(1 * time.Second)
	deleteBtns, err := wd.FindElements(selenium.ByClassName, "btn-danger")
	if err != nil || len(deleteBtns) == 0 {
		t.Fatalf("Failed to find delete buttons: %v", err)
	}
	for _, btn := range deleteBtns {
		text, _ := btn.Text()
		if text == "Delete" {
			btn.Click()
			break
		}
	}

	// Confirm delete in modal
	time.Sleep(500 * time.Millisecond)
	confirmBtn := waitForElement(t, selenium.ByID, "confirm-yes", 2*time.Second)
	confirmBtn.Click()

	waitForNotificationWithText(t, "Book deleted successfully", 5*time.Second)
}

// Test 5: Create Reader
func TestCreateReader(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	readersTab := waitForElement(t, selenium.ByXPATH, "//button[contains(text(), 'Readers')]", 5*time.Second)
	readersTab.Click()
	time.Sleep(1 * time.Second)

	addReaderBtn := waitForElement(t, selenium.ByID, "add-reader-btn", 5*time.Second)
	addReaderBtn.Click()

	nameInput := waitForElement(t, selenium.ByID, "reader-name", 2*time.Second)
	nameInput.SendKeys("John")

	surnameInput := waitForElement(t, selenium.ByID, "reader-surname", 1*time.Second)
	surnameInput.SendKeys("Doe")

	submitBtn, err := wd.FindElement(selenium.ByID, "reader-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find reader submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Reader created successfully", 5*time.Second)
}

// Test 6: Add Book to Reader's Currently Reading List
func TestAddBookToReaderReadingList(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	// Create a book first
	addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
	addBookBtn.Click()

	// Step 1: Fill title
	titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.SendKeys("Book for Reading List")

	// Navigate through steps to submit
	nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	submitBtn, err := wd.FindElement(selenium.ByID, "book-form-submit")
	if err != nil {
		t.Fatalf("Failed to find book submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Book created successfully", 5*time.Second)
	time.Sleep(1 * time.Second)

	// Switch to readers tab
	readersTab := waitForElement(t, selenium.ByXPATH, "//button[contains(text(), 'Readers')]", 5*time.Second)
	readersTab.Click()
	time.Sleep(1 * time.Second)

	// Create a reader
	addReaderBtn := waitForElement(t, selenium.ByID, "add-reader-btn", 5*time.Second)
	addReaderBtn.Click()

	nameInput := waitForElement(t, selenium.ByID, "reader-name", 2*time.Second)
	nameInput.SendKeys("Jane")

	surnameInput := waitForElement(t, selenium.ByID, "reader-surname", 1*time.Second)
	surnameInput.SendKeys("Smith")

	submitBtn, err = wd.FindElement(selenium.ByID, "reader-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find reader submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Reader created successfully", 5*time.Second)
	time.Sleep(1 * time.Second)

	// Click "Add Book" button
	addBookToReaderBtns, err := wd.FindElements(selenium.ByClassName, "btn-secondary")
	if err != nil {
		t.Fatalf("Failed to find add book buttons: %v", err)
	}
	var addBookToReaderBtn selenium.WebElement
	for _, btn := range addBookToReaderBtns {
		text, _ := btn.Text()
		if contains(text, "Add Book") {
			addBookToReaderBtn = btn
			break
		}
	}

	if addBookToReaderBtn == nil {
		t.Fatal("Add Book button not found")
	}

	addBookToReaderBtn.Click()
	time.Sleep(1 * time.Second)

	// Select a book from the modal
	bookOptions, err := wd.FindElements(selenium.ByClassName, "book-option")
	if err != nil || len(bookOptions) == 0 {
		t.Fatal("No book options found in modal")
	}

	bookOptions[0].Click()

	waitForNotificationWithText(t, "Book added to reading list!", 5*time.Second)
}

// Test 7: Remove Book from Reader's Reading List
func TestRemoveBookFromReaderReadingList(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	// Create a book
	addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
	addBookBtn.Click()

	// Step 1: Fill title
	titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.SendKeys("Book to Remove")

	// Navigate through steps to submit
	nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	nextBtn.Click()
	time.Sleep(300 * time.Millisecond)

	submitBtn, err := wd.FindElement(selenium.ByID, "book-form-submit")
	if err != nil {
		t.Fatalf("Failed to find book submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Book created successfully", 5*time.Second)
	time.Sleep(1 * time.Second)

	// Switch to readers tab
	readersTab := waitForElement(t, selenium.ByXPATH, "//button[contains(text(), 'Readers')]", 5*time.Second)
	readersTab.Click()
	time.Sleep(1 * time.Second)

	// Create a reader
	addReaderBtn := waitForElement(t, selenium.ByID, "add-reader-btn", 5*time.Second)
	addReaderBtn.Click()

	nameInput := waitForElement(t, selenium.ByID, "reader-name", 2*time.Second)
	nameInput.SendKeys("Bob")

	surnameInput := waitForElement(t, selenium.ByID, "reader-surname", 1*time.Second)
	surnameInput.SendKeys("Jones")

	submitBtn, err = wd.FindElement(selenium.ByID, "reader-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find reader submit button: %v", err)
	}
	submitBtn.Click()

	waitForNotificationWithText(t, "Reader created successfully", 5*time.Second)
	time.Sleep(1 * time.Second)

	// Add book to reader
	addBookToReaderBtns, err := wd.FindElements(selenium.ByClassName, "btn-secondary")
	if err != nil {
		t.Fatalf("Failed to find add book buttons: %v", err)
	}
	var addBookToReaderBtn selenium.WebElement
	for _, btn := range addBookToReaderBtns {
		text, _ := btn.Text()
		if contains(text, "Add Book") {
			addBookToReaderBtn = btn
			break
		}
	}

	if addBookToReaderBtn == nil {
		t.Fatal("Add Book button not found")
	}

	addBookToReaderBtn.Click()
	time.Sleep(1 * time.Second)

	bookOptions, err := wd.FindElements(selenium.ByClassName, "book-option")
	if err != nil || len(bookOptions) == 0 {
		t.Fatal("No book options found in modal")
	}
	bookOptions[0].Click()

	waitForNotificationWithText(t, "Book added to reading list!", 5*time.Second)
	time.Sleep(1 * time.Second)

	// Remove the book
	removeBtn := waitForElement(t, selenium.ByClassName, "btn-remove-book", 3*time.Second)
	removeBtn.Click()

	// Confirm removal in modal
	time.Sleep(500 * time.Millisecond)
	confirmBtn := waitForElement(t, selenium.ByID, "confirm-yes", 2*time.Second)
	confirmBtn.Click()

	waitForNotificationWithText(t, "Book removed from reading list!", 5*time.Second)
}

// Test 8: Search Books
func TestSearchBooks(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	// Create two books
	for i := 1; i <= 2; i++ {
		addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
		addBookBtn.Click()

		// Step 1: Fill title
		titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
		if i == 1 {
			titleInput.SendKeys("Searchable Book")
		} else {
			titleInput.SendKeys("Another Book")
		}

		// Navigate through steps to submit
		nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
		nextBtn.Click()
		time.Sleep(300 * time.Millisecond)

		nextBtn.Click()
		time.Sleep(300 * time.Millisecond)

		submitBtn, err := wd.FindElement(selenium.ByID, "book-form-submit")
		if err != nil {
			t.Fatalf("Failed to find book submit button: %v", err)
		}
		submitBtn.Click()

		waitForNotificationWithText(t, "Book created successfully", 5*time.Second)
		time.Sleep(1 * time.Second)
	}

	// Search for "Searchable" - search is now immediate
	searchInput := waitForElement(t, selenium.ByID, "book-search-query", 3*time.Second)
	searchInput.Clear()
	time.Sleep(300 * time.Millisecond)
	searchInput.SendKeys("Searchable")

	// Wait for immediate filter to apply
	time.Sleep(2 * time.Second)

	// Verify only one book is visible
	bookCards, err := wd.FindElements(selenium.ByClassName, "item-card")
	if err != nil {
		t.Fatalf("Failed to find book cards: %v", err)
	}
	if len(bookCards) < 1 {
		t.Fatalf("Expected at least 1 book card, got 0")
	}
}

// Test 9: User Profile
func TestUserProfile(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	userMenuBtn := waitForElement(t, selenium.ByID, "user-menu-btn", 5*time.Second)
	userMenuBtn.Click()
	time.Sleep(500 * time.Millisecond)

	profileBtn := waitForElement(t, selenium.ByID, "profile-menu-item", 2*time.Second)
	profileBtn.Click()
	time.Sleep(500 * time.Millisecond)

	profileModal := waitForElement(t, selenium.ByID, "profile-modal", 2*time.Second)
	displayed, _ := profileModal.IsDisplayed()
	if !displayed {
		t.Fatal("Profile modal should be visible")
	}
}

// Test 10: Negative - Create Book with Empty Title
func TestCreateBookNegativeEmptyTitle(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	loginHelper(t)

	addBookBtn := waitForElement(t, selenium.ByID, "add-book-btn", 5*time.Second)
	addBookBtn.Click()

	titleInput := waitForElement(t, selenium.ByID, "book-title", 2*time.Second)
	titleInput.Clear()

	nextBtn := waitForElement(t, selenium.ByID, "book-form-next", 1*time.Second)
	nextBtn.Click()

	time.Sleep(500 * time.Millisecond)

	validationMsg, _ := titleInput.GetAttribute("validationMessage")
	if validationMsg == "" {
		t.Fatal("Validation message should be displayed for empty title")
	}
}

// Test 11: Negative - Invalid Login Credentials
func TestLoginNegativeInvalidCredentials(t *testing.T) {
	setupWebDriver(t)
	defer teardownWebDriver()

	if err := wd.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	loginTab := waitForElement(t, selenium.ByXPATH, "//button[@data-auth-tab='login']", 5*time.Second)
	loginTab.Click()
	time.Sleep(500 * time.Millisecond)

	loginUsername := waitForElement(t, selenium.ByID, "login-username", 2*time.Second)
	loginUsername.SendKeys("nonexistentuser")

	loginPassword := waitForElement(t, selenium.ByID, "login-password", 1*time.Second)
	loginPassword.SendKeys("wrongpassword")

	loginSubmit, err := wd.FindElement(selenium.ByID, "login-submit-btn")
	if err != nil {
		t.Fatalf("Failed to find login submit button: %v", err)
	}
	loginSubmit.Click()

	time.Sleep(2 * time.Second)

	notification := waitForElement(t, selenium.ByID, "notification", 5*time.Second)
	className, _ := notification.GetAttribute("class")
	if !contains(className, "error") {
		t.Fatal("Error notification should be displayed for invalid credentials")
	}
}
