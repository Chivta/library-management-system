const API_BASE = 'http://localhost:8080';

let currentTab = 'books';
let editingBookId = null;
let editingReaderId = null;

function showNotification(message, type = 'info') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification ${type}`;
    notification.classList.remove('hidden');

    setTimeout(() => {
        notification.classList.add('hidden');
    }, 5000);
}

function showLoading(element) {
    element.classList.remove('hidden');
}

function hideLoading(element) {
    element.classList.add('hidden');
}

function clearFormErrors(formPrefix) {
    const errorElements = document.querySelectorAll(`[id^="${formPrefix}-"][id$="-error"]`);
    errorElements.forEach(el => el.textContent = '');
}

function displayValidationErrors(errors, formPrefix) {
    clearFormErrors(formPrefix);

    if (Array.isArray(errors)) {
        errors.forEach(error => {
            const fieldName = error.field.toLowerCase();
            const errorElement = document.getElementById(`${formPrefix}-${fieldName}-error`);
            if (errorElement) {
                errorElement.textContent = error.message;
            }
        });
    }
}

async function apiRequest(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });

        if (response.status === 204) {
            return { success: true };
        }

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 403) {
                throw new Error(data.error || 'This endpoint is currently disabled');
            } else if (response.status === 404) {
                throw new Error(data.error || 'Resource not found');
            } else if (response.status === 400) {
                if (data.errors) {
                    throw { validationErrors: data.errors, message: 'Validation failed' };
                }
                throw new Error(data.error || 'Invalid request');
            } else if (response.status === 500) {
                throw new Error(data.error || 'Server error occurred');
            } else {
                throw new Error(data.error || 'An error occurred');
            }
        }

        return data;
    } catch (error) {
        if (error.validationErrors) {
            throw error;
        }
        if (error.message === 'Failed to fetch') {
            throw new Error('Unable to connect to the server. Please ensure the server is running.');
        }
        throw error;
    }
}

async function getAllBooks() {
    return await apiRequest('/books/');
}

async function getBookById(id) {
    return await apiRequest(`/books/${id}`);
}

async function createBook(bookData) {
    return await apiRequest('/books/', {
        method: 'POST',
        body: JSON.stringify(bookData)
    });
}

async function updateBook(id, bookData) {
    return await apiRequest(`/books/${id}`, {
        method: 'PUT',
        body: JSON.stringify(bookData)
    });
}

async function deleteBook(id) {
    return await apiRequest(`/books/${id}`, {
        method: 'DELETE'
    });
}

async function deleteAllBooks() {
    return await apiRequest('/books/', {
        method: 'DELETE'
    });
}

async function getAllReaders() {
    return await apiRequest('/readers/');
}

async function getReaderById(id) {
    return await apiRequest(`/readers/${id}`);
}

async function createReader(readerData) {
    return await apiRequest('/readers/', {
        method: 'POST',
        body: JSON.stringify(readerData)
    });
}

async function updateReader(id, readerData) {
    return await apiRequest(`/readers/${id}`, {
        method: 'PUT',
        body: JSON.stringify(readerData)
    });
}

async function deleteReader(id) {
    return await apiRequest(`/readers/${id}`, {
        method: 'DELETE'
    });
}

async function deleteAllReaders() {
    return await apiRequest('/readers/', {
        method: 'DELETE'
    });
}

function renderBooks(books) {
    const container = document.getElementById('books-list');

    if (!books || books.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <h3>No Books Found</h3>
                <p>Start by adding your first book!</p>
            </div>
        `;
        return;
    }

    const html = `
        <div class="items-grid">
            ${books.map(book => `
                <div class="item-card">
                    <div class="item-card-header">
                        <span class="item-id">#${book.id}</span>
                    </div>
                    <h4>${escapeHtml(book.title)}</h4>
                    <p>${escapeHtml(book.description || 'No description')}</p>
                    <div class="item-card-actions">
                        <button class="btn btn-primary btn-small" onclick="editBook(${book.id})">Edit</button>
                        <button class="btn btn-danger btn-small" onclick="confirmDeleteBook(${book.id})">Delete</button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;

    container.innerHTML = html;
}

function renderReaders(readers) {
    const container = document.getElementById('readers-list');

    if (!readers || readers.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <h3>No Readers Found</h3>
                <p>Start by adding your first reader!</p>
            </div>
        `;
        return;
    }

    const html = `
        <div class="items-grid">
            ${readers.map(reader => `
                <div class="item-card">
                    <div class="item-card-header">
                        <span class="item-id">#${reader.id}</span>
                    </div>
                    <h4>${escapeHtml(reader.name)} ${escapeHtml(reader.surname)}</h4>
                    <div class="item-card-actions">
                        <button class="btn btn-primary btn-small" onclick="editReader(${reader.id})">Edit</button>
                        <button class="btn btn-danger btn-small" onclick="confirmDeleteReader(${reader.id})">Delete</button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;

    container.innerHTML = html;
}

function renderBookDetail(book) {
    const container = document.getElementById('book-detail');
    container.innerHTML = `
        <div class="detail-item"><strong>ID:</strong> ${book.id}</div>
        <div class="detail-item"><strong>Title:</strong> ${escapeHtml(book.title)}</div>
        <div class="detail-item"><strong>Description:</strong> ${escapeHtml(book.description || 'No description')}</div>
        <div class="detail-actions">
            <button class="btn btn-primary btn-small" onclick="editBook(${book.id})">Edit</button>
            <button class="btn btn-danger btn-small" onclick="confirmDeleteBook(${book.id})">Delete</button>
        </div>
    `;
    container.classList.remove('hidden');
}

function renderReaderDetail(reader) {
    const container = document.getElementById('reader-detail');
    container.innerHTML = `
        <div class="detail-item"><strong>ID:</strong> ${reader.id}</div>
        <div class="detail-item"><strong>Name:</strong> ${escapeHtml(reader.name)}</div>
        <div class="detail-item"><strong>Surname:</strong> ${escapeHtml(reader.surname)}</div>
        <div class="detail-actions">
            <button class="btn btn-primary btn-small" onclick="editReader(${reader.id})">Edit</button>
            <button class="btn btn-danger btn-small" onclick="confirmDeleteReader(${reader.id})">Delete</button>
        </div>
    `;
    container.classList.remove('hidden');
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function loadBooks() {
    const loading = document.getElementById('books-loading');
    showLoading(loading);

    try {
        const books = await getAllBooks();
        renderBooks(books);
    } catch (error) {
        showNotification(error.message, 'error');
        renderBooks([]);
    } finally {
        hideLoading(loading);
    }
}

async function searchBook() {
    const idInput = document.getElementById('search-book-id');
    const id = idInput.value.trim();

    if (!id) {
        showNotification('Please enter a book ID', 'error');
        return;
    }

    try {
        const book = await getBookById(id);
        renderBookDetail(book);
        showNotification('Book found successfully', 'success');
    } catch (error) {
        document.getElementById('book-detail').classList.add('hidden');
        showNotification(error.message, 'error');
    }
}

function showBookForm(isEdit = false) {
    const container = document.getElementById('book-form-container');
    const title = document.getElementById('book-form-title');
    const form = document.getElementById('book-form');

    title.textContent = isEdit ? 'Edit Book' : 'Add New Book';
    form.reset();
    clearFormErrors('book');
    container.classList.remove('hidden');

    if (!isEdit) {
        document.getElementById('book-id').value = '';
        editingBookId = null;
    }
}

function hideBookForm() {
    document.getElementById('book-form-container').classList.add('hidden');
    editingBookId = null;
}

async function editBook(id) {
    try {
        const book = await getBookById(id);
        showBookForm(true);

        document.getElementById('book-id').value = book.id;
        document.getElementById('book-title').value = book.title;
        document.getElementById('book-description').value = book.description || '';

        editingBookId = id;
    } catch (error) {
        showNotification(error.message, 'error');
    }
}

async function handleBookFormSubmit(e) {
    e.preventDefault();

    const submitBtn = document.getElementById('book-submit-btn');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Saving...';

    const bookData = {
        title: document.getElementById('book-title').value.trim(),
        description: document.getElementById('book-description').value.trim()
    };

    try {
        if (editingBookId) {
            await updateBook(editingBookId, bookData);
            showNotification('Book updated successfully', 'success');
        } else {
            await createBook(bookData);
            showNotification('Book created successfully', 'success');
        }

        hideBookForm();
        await loadBooks();
    } catch (error) {
        if (error.validationErrors) {
            displayValidationErrors(error.validationErrors, 'book');
            showNotification('Please fix the validation errors', 'error');
        } else {
            showNotification(error.message, 'error');
        }
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Save Book';
    }
}

async function confirmDeleteBook(id) {
    showConfirmModal(
        'Are you sure you want to delete this book? This action cannot be undone.',
        async () => {
            try {
                await deleteBook(id);
                showNotification('Book deleted successfully', 'success');
                await loadBooks();
                document.getElementById('book-detail').classList.add('hidden');
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    );
}

async function confirmDeleteAllBooks() {
    showConfirmModal(
        'Are you sure you want to delete ALL books? This action cannot be undone and will remove all books from the database.',
        async () => {
            try {
                await deleteAllBooks();
                showNotification('All books deleted successfully', 'success');
                await loadBooks();
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    );
}

async function loadReaders() {
    const loading = document.getElementById('readers-loading');
    showLoading(loading);

    try {
        const readers = await getAllReaders();
        renderReaders(readers);
    } catch (error) {
        showNotification(error.message, 'error');
        renderReaders([]);
    } finally {
        hideLoading(loading);
    }
}

async function searchReader() {
    const idInput = document.getElementById('search-reader-id');
    const id = idInput.value.trim();

    if (!id) {
        showNotification('Please enter a reader ID', 'error');
        return;
    }

    try {
        const reader = await getReaderById(id);
        renderReaderDetail(reader);
        showNotification('Reader found successfully', 'success');
    } catch (error) {
        document.getElementById('reader-detail').classList.add('hidden');
        showNotification(error.message, 'error');
    }
}

function showReaderForm(isEdit = false) {
    const container = document.getElementById('reader-form-container');
    const title = document.getElementById('reader-form-title');
    const form = document.getElementById('reader-form');

    title.textContent = isEdit ? 'Edit Reader' : 'Add New Reader';
    form.reset();
    clearFormErrors('reader');
    container.classList.remove('hidden');

    if (!isEdit) {
        document.getElementById('reader-id').value = '';
        editingReaderId = null;
    }
}

function hideReaderForm() {
    document.getElementById('reader-form-container').classList.add('hidden');
    editingReaderId = null;
}

async function editReader(id) {
    try {
        const reader = await getReaderById(id);
        showReaderForm(true);

        document.getElementById('reader-id').value = reader.id;
        document.getElementById('reader-name').value = reader.name;
        document.getElementById('reader-surname').value = reader.surname;

        editingReaderId = id;
    } catch (error) {
        showNotification(error.message, 'error');
    }
}

async function handleReaderFormSubmit(e) {
    e.preventDefault();

    const submitBtn = document.getElementById('reader-submit-btn');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Saving...';

    const readerData = {
        name: document.getElementById('reader-name').value.trim(),
        surname: document.getElementById('reader-surname').value.trim()
    };

    try {
        if (editingReaderId) {
            await updateReader(editingReaderId, readerData);
            showNotification('Reader updated successfully', 'success');
        } else {
            await createReader(readerData);
            showNotification('Reader created successfully', 'success');
        }

        hideReaderForm();
        await loadReaders();
    } catch (error) {
        if (error.validationErrors) {
            displayValidationErrors(error.validationErrors, 'reader');
            showNotification('Please fix the validation errors', 'error');
        } else {
            showNotification(error.message, 'error');
        }
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Save Reader';
    }
}

async function confirmDeleteReader(id) {
    showConfirmModal(
        'Are you sure you want to delete this reader? This action cannot be undone.',
        async () => {
            try {
                await deleteReader(id);
                showNotification('Reader deleted successfully', 'success');
                await loadReaders();
                document.getElementById('reader-detail').classList.add('hidden');
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    );
}

async function confirmDeleteAllReaders() {
    showConfirmModal(
        'Are you sure you want to delete ALL readers? This action cannot be undone and will remove all readers from the database.',
        async () => {
            try {
                await deleteAllReaders();
                showNotification('All readers deleted successfully', 'success');
                await loadReaders();
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    );
}

function showConfirmModal(message, onConfirm) {
    const modal = document.getElementById('confirm-modal');
    const messageEl = document.getElementById('confirm-message');

    messageEl.textContent = message;
    modal.classList.remove('hidden');

    const handleConfirm = async () => {
        modal.classList.add('hidden');
        await onConfirm();
        cleanup();
    };

    const handleCancel = () => {
        modal.classList.add('hidden');
        cleanup();
    };

    const cleanup = () => {
        document.getElementById('confirm-yes').removeEventListener('click', handleConfirm);
        document.getElementById('confirm-no').removeEventListener('click', handleCancel);
    };

    document.getElementById('confirm-yes').addEventListener('click', handleConfirm);
    document.getElementById('confirm-no').addEventListener('click', handleCancel);
}

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

    currentTab = tabName;

    if (tabName === 'books') {
        loadBooks();
    } else if (tabName === 'readers') {
        loadReaders();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });

    document.getElementById('add-book-btn').addEventListener('click', () => showBookForm(false));
    document.getElementById('refresh-books-btn').addEventListener('click', loadBooks);
    document.getElementById('delete-all-books-btn').addEventListener('click', confirmDeleteAllBooks);
    document.getElementById('book-form').addEventListener('submit', handleBookFormSubmit);
    document.getElementById('book-cancel-btn').addEventListener('click', hideBookForm);
    document.getElementById('search-book-btn').addEventListener('click', searchBook);
    document.getElementById('search-book-id').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') searchBook();
    });

    document.getElementById('add-reader-btn').addEventListener('click', () => showReaderForm(false));
    document.getElementById('refresh-readers-btn').addEventListener('click', loadReaders);
    document.getElementById('delete-all-readers-btn').addEventListener('click', confirmDeleteAllReaders);
    document.getElementById('reader-form').addEventListener('submit', handleReaderFormSubmit);
    document.getElementById('reader-cancel-btn').addEventListener('click', hideReaderForm);
    document.getElementById('search-reader-btn').addEventListener('click', searchReader);
    document.getElementById('search-reader-id').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') searchReader();
    });

    loadBooks();
});
