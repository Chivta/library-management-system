function switchTab(tabName) {
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.tab === tabName) {
            btn.classList.add('active');
        }
    });

    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });
    document.getElementById(`${tabName}-section`).classList.add('active');

    AppState.currentTab = tabName;

    if (tabName === 'books') {
        Books.load();
    } else if (tabName === 'readers') {
        Readers.load();
    } else if (tabName === 'statistics') {
        Statistics.load();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.body.classList.add('loaded');

    document.querySelectorAll('.auth-tab').forEach(btn => {
        btn.addEventListener('click', () => switchAuthTab(btn.dataset.authTab));
    });

    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');

    if (loginForm) loginForm.addEventListener('submit', handleLoginSubmit);
    if (registerForm) registerForm.addEventListener('submit', handleRegisterSubmit);

    const userMenuBtn = document.getElementById('user-menu-btn');
    const userDropdown = document.getElementById('user-dropdown');
    const profileMenuItem = document.getElementById('profile-menu-item');
    const logoutMenuItem = document.getElementById('logout-menu-item');
    const profileModalClose = document.getElementById('profile-modal-close');

    if (userMenuBtn && userDropdown) {
        userMenuBtn.addEventListener('click', () => {
            userDropdown.classList.toggle('hidden');
        });
    }

    if (profileMenuItem && userDropdown) {
        profileMenuItem.addEventListener('click', () => {
            userDropdown.classList.add('hidden');
            UI.showProfileModal();
        });
    }

    if (logoutMenuItem && userDropdown) {
        logoutMenuItem.addEventListener('click', () => {
            userDropdown.classList.add('hidden');
            Auth.logout();
        });
    }

    if (profileModalClose) profileModalClose.addEventListener('click', UI.hideProfileModal);

    document.addEventListener('click', (e) => {
        const dropdown = document.getElementById('user-dropdown');
        if (dropdown && !e.target.closest('.user-menu')) {
            dropdown.classList.add('hidden');
        }
    });

    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });

    const addBookBtn = document.getElementById('add-book-btn');
    const refreshBooksBtn = document.getElementById('refresh-books-btn');
    const exportBooksBtn = document.getElementById('export-books-btn');
    const deleteAllBooksBtn = document.getElementById('delete-all-books-btn');
    const bookForm = document.getElementById('book-form');
    const bookFormCancel = document.getElementById('book-form-cancel');
    const bookModalClose = document.getElementById('book-modal-close');
    const bookFormNext = document.getElementById('book-form-next');
    const bookFormPrev = document.getElementById('book-form-prev');

    if (addBookBtn) addBookBtn.addEventListener('click', () => Books.showModal(false));
    if (refreshBooksBtn) refreshBooksBtn.addEventListener('click', () => Books.load());
    if (exportBooksBtn) exportBooksBtn.addEventListener('click', () => Books.export());
    if (deleteAllBooksBtn) deleteAllBooksBtn.addEventListener('click', () => Books.confirmDeleteAll());
    if (bookForm) bookForm.addEventListener('submit', (e) => Books.handleFormSubmit(e));
    if (bookFormCancel) bookFormCancel.addEventListener('click', () => Books.hideModal());
    if (bookModalClose) bookModalClose.addEventListener('click', () => Books.hideModal());
    if (bookFormNext) bookFormNext.addEventListener('click', () => Books.nextFormStep());
    if (bookFormPrev) bookFormPrev.addEventListener('click', () => Books.prevFormStep());

    // Immediate book search and filter
    const bookSearchInput = document.getElementById('book-search-query');
    const bookSortBy = document.getElementById('book-sort-by');
    const bookItemsPerPage = document.getElementById('book-items-per-page');

    if (bookSearchInput) {
        bookSearchInput.addEventListener('input', () => {
            AppState.bookFilters.search = bookSearchInput.value;
            Books.applyFilters();
        });
    }

    if (bookSortBy) {
        bookSortBy.addEventListener('change', () => {
            AppState.bookFilters.sortBy = bookSortBy.value;
            Books.applyFilters();
        });
    }

    if (bookItemsPerPage) {
        bookItemsPerPage.addEventListener('change', () => {
            AppState.booksItemsPerPage = parseInt(bookItemsPerPage.value);
            Books.applyFilters();
        });
    }

    const addReaderBtn = document.getElementById('add-reader-btn');
    const refreshReadersBtn = document.getElementById('refresh-readers-btn');
    const exportReadersBtn = document.getElementById('export-readers-btn');
    const deleteAllReadersBtn = document.getElementById('delete-all-readers-btn');
    const readerForm = document.getElementById('reader-form');
    const readerCancelBtn = document.getElementById('reader-cancel-btn');
    const readerModalClose = document.getElementById('reader-modal-close');

    if (addReaderBtn) addReaderBtn.addEventListener('click', () => Readers.showModal(false));
    if (refreshReadersBtn) refreshReadersBtn.addEventListener('click', () => Readers.load());
    if (exportReadersBtn) exportReadersBtn.addEventListener('click', () => Readers.export());
    if (deleteAllReadersBtn) deleteAllReadersBtn.addEventListener('click', () => Readers.confirmDeleteAll());
    if (readerForm) readerForm.addEventListener('submit', (e) => Readers.handleFormSubmit(e));
    if (readerCancelBtn) readerCancelBtn.addEventListener('click', () => Readers.hideModal());
    if (readerModalClose) readerModalClose.addEventListener('click', () => Readers.hideModal());

    // Immediate reader search and filter
    const readerSearchInput = document.getElementById('reader-search-query');
    const readerSortBy = document.getElementById('reader-sort-by');
    const readerItemsPerPage = document.getElementById('reader-items-per-page');

    if (readerSearchInput) {
        readerSearchInput.addEventListener('input', () => {
            AppState.readerFilters.search = readerSearchInput.value;
            Readers.applyFilters();
        });
    }

    if (readerSortBy) {
        readerSortBy.addEventListener('change', () => {
            AppState.readerFilters.sortBy = readerSortBy.value;
            Readers.applyFilters();
        });
    }

    if (readerItemsPerPage) {
        readerItemsPerPage.addEventListener('change', () => {
            AppState.readersItemsPerPage = parseInt(readerItemsPerPage.value);
            Readers.applyFilters();
        });
    }

    const refreshStatsBtn = document.getElementById('refresh-stats-btn');
    if (refreshStatsBtn) refreshStatsBtn.addEventListener('click', () => Statistics.load());

    if (Auth.isAuthenticated()) {
        Auth.getProfile().then(profile => {
            AppState.currentUser = profile;
            document.getElementById('username-display').textContent = profile.username;
            UI.showAppContainer();
            Books.load();
        }).catch(() => {
            Auth.clearToken();
            UI.showAuthContainer();
        });
    } else {
        UI.showAuthContainer();
    }
});
