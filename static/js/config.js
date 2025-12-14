const API_BASE = 'http://localhost:8080';

const AppState = {
    authToken: null,
    currentUser: null,
    currentTab: 'books',

    allBooks: [],
    filteredBooks: [],
    booksCurrentPage: 1,
    booksItemsPerPage: 10,
    bookFilters: {
        search: '',
        sortBy: 'id-desc'
    },
    editingBookId: null,
    currentFormStep: 1,

    allReaders: [],
    filteredReaders: [],
    readersCurrentPage: 1,
    readersItemsPerPage: 10,
    readerFilters: {
        search: '',
        sortBy: 'id-asc'
    },
    editingReaderId: null
};
