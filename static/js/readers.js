const Readers = {
    applyFilters() {
        const search = AppState.readerFilters.search.toLowerCase();

        AppState.filteredReaders = AppState.allReaders.filter(reader => {
            const matchesSearch = !search ||
                reader.name.toLowerCase().includes(search) ||
                reader.surname.toLowerCase().includes(search);
            return matchesSearch;
        });

        const [field, direction] = AppState.readerFilters.sortBy.split('-');
        AppState.filteredReaders.sort((a, b) => {
            let aVal, bVal;

            if (field === 'id') {
                aVal = a.id;
                bVal = b.id;
            } else if (field === 'name') {
                aVal = a.name.toLowerCase();
                bVal = b.name.toLowerCase();
            }

            if (direction === 'asc') {
                return aVal > bVal ? 1 : -1;
            } else {
                return aVal < bVal ? 1 : -1;
            }
        });

        AppState.readersCurrentPage = 1;
        this.render();
    },

    render() {
        const container = document.getElementById('readers-list');

        if (!AppState.filteredReaders || AppState.filteredReaders.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3>No Readers Found</h3>
                    <p>Start by adding your first reader or adjust your filters!</p>
                </div>
            `;
            document.getElementById('readers-pagination').innerHTML = '';
            return;
        }

        // Pagination logic
        const startIndex = (AppState.readersCurrentPage - 1) * AppState.readersItemsPerPage;
        const endIndex = startIndex + AppState.readersItemsPerPage;
        const paginatedReaders = AppState.filteredReaders.slice(startIndex, endIndex);

        const html = `
            <div class="items-grid">
                ${paginatedReaders.map(reader => `
                    <div class="item-card" data-reader-id="${reader.id}">
                        <div class="item-card-header">
                            <span class="item-id">#${reader.id}</span>
                        </div>
                        <h4>${escapeHtml(reader.name)} ${escapeHtml(reader.surname)}</h4>
                        <div class="currently-reading">
                            <strong>Currently Reading:</strong>
                            ${reader.currently_reading && reader.currently_reading.length > 0 ? `
                                <ul class="reading-list">
                                    ${reader.currently_reading.map(book => `
                                        <li>
                                            ${escapeHtml(book.title)}
                                            <button class="btn-remove-book" onclick="Readers.removeBook(${reader.id}, ${book.id})" title="Remove">×</button>
                                        </li>
                                    `).join('')}
                                </ul>
                            ` : '<span class="empty-reading">No books</span>'}
                            <button class="btn btn-secondary btn-small" onclick="Readers.showAddBookModal(${reader.id})">+ Add Book</button>
                        </div>
                        <div class="item-card-actions">
                            <button class="btn btn-primary btn-small" onclick="Readers.edit(${reader.id})">Edit</button>
                            <button class="btn btn-danger btn-small" onclick="Readers.confirmDelete(${reader.id})">Delete</button>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;

        container.innerHTML = html;
        this.renderPagination();
    },

    renderPagination() {
        const container = document.getElementById('readers-pagination');
        const totalPages = Math.ceil(AppState.filteredReaders.length / AppState.readersItemsPerPage);

        if (totalPages <= 1) {
            container.innerHTML = '';
            return;
        }

        let html = '<div class="pagination-controls">';
        html += `<button class="btn btn-small" onclick="Readers.changePage(${AppState.readersCurrentPage - 1})" ${AppState.readersCurrentPage === 1 ? 'disabled' : ''}>← Previous</button>`;
        html += '<div class="pagination-pages">';

        for (let i = 1; i <= totalPages; i++) {
            if (i === 1 || i === totalPages || (i >= AppState.readersCurrentPage - 2 && i <= AppState.readersCurrentPage + 2)) {
                html += `<button class="btn btn-small ${i === AppState.readersCurrentPage ? 'btn-primary' : ''}" onclick="Readers.changePage(${i})">${i}</button>`;
            } else if (i === AppState.readersCurrentPage - 3 || i === AppState.readersCurrentPage + 3) {
                html += '<span>...</span>';
            }
        }

        html += '</div>';
        html += `<button class="btn btn-small" onclick="Readers.changePage(${AppState.readersCurrentPage + 1})" ${AppState.readersCurrentPage === totalPages ? 'disabled' : ''}>Next →</button>`;
        html += '</div>';

        container.innerHTML = html;
    },

    changePage(page) {
        const totalPages = Math.ceil(AppState.filteredReaders.length / AppState.readersItemsPerPage);
        if (page >= 1 && page <= totalPages) {
            AppState.readersCurrentPage = page;
            this.render();
            document.getElementById('readers-list').scrollIntoView({ behavior: 'smooth' });
        }
    },

    async load() {
        const loading = document.getElementById('readers-loading');
        UI.showLoading(loading);

        try {
            AppState.allReaders = await ReadersAPI.getAll();
            AppState.filteredReaders = [...AppState.allReaders];
            this.applyFilters();
        } catch (error) {
            UI.showNotification(error.message, 'error');
            AppState.allReaders = [];
            AppState.filteredReaders = [];
            this.render();
        } finally {
            UI.hideLoading(loading);
        }
    },

    showModal(isEdit = false) {
        const modal = document.getElementById('reader-modal');
        const title = document.getElementById('reader-modal-title');

        title.textContent = isEdit ? 'Edit Reader' : 'Add New Reader';
        UI.clearFormErrors('reader');

        if (!isEdit) {
            document.getElementById('reader-form').reset();
            AppState.editingReaderId = null;
        }

        modal.classList.remove('hidden');
    },

    hideModal() {
        document.getElementById('reader-modal').classList.add('hidden');
        AppState.editingReaderId = null;
    },

    async edit(id) {
        try {
            const reader = await ReadersAPI.getById(id);
            this.showModal(true);

            document.getElementById('reader-id').value = reader.id;
            document.getElementById('reader-name').value = reader.name;
            document.getElementById('reader-surname').value = reader.surname;

            AppState.editingReaderId = id;
        } catch (error) {
            UI.showNotification(error.message, 'error');
        }
    },

    async handleFormSubmit(e) {
        e.preventDefault();

        const submitBtn = document.getElementById('reader-submit-btn');
        submitBtn.disabled = true;
        submitBtn.textContent = 'Saving...';

        const readerData = {
            name: document.getElementById('reader-name').value.trim(),
            surname: document.getElementById('reader-surname').value.trim()
        };

        try {
            if (AppState.editingReaderId) {
                await ReadersAPI.update(AppState.editingReaderId, readerData);
                UI.showNotification('Reader updated successfully', 'success');
            } else {
                await ReadersAPI.create(readerData);
                UI.showNotification('Reader created successfully', 'success');
            }

            this.hideModal();
            await this.load();
        } catch (error) {
            if (error.validationErrors) {
                UI.displayValidationErrors(error.validationErrors, 'reader');
                UI.showNotification('Please fix the validation errors', 'error');
            } else {
                UI.showNotification(error.message, 'error');
            }
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Save Reader';
        }
    },

    confirmDelete(id) {
        UI.showConfirmModal(
            'Are you sure you want to delete this reader? This action cannot be undone.',
            async () => {
                try {
                    await ReadersAPI.delete(id);
                    UI.showNotification('Reader deleted successfully', 'success');
                    await this.load();
                } catch (error) {
                    UI.showNotification(error.message, 'error');
                }
            }
        );
    },

    confirmDeleteAll() {
        UI.showConfirmModal(
            'Are you sure you want to delete ALL readers? This action cannot be undone and will remove all readers from the database.',
            async () => {
                try {
                    await ReadersAPI.deleteAll();
                    UI.showNotification('All readers deleted successfully', 'success');
                    await this.load();
                } catch (error) {
                    UI.showNotification(error.message, 'error');
                }
            }
        );
    },

    export() {
        const csv = convertReadersToCSV(AppState.filteredReaders);
        downloadCSV(csv, 'readers.csv');
        UI.showNotification('Readers exported successfully', 'success');
    },

    async showAddBookModal(readerID) {
        AppState.editingReaderId = readerID;

        // Get all books
        try {
            const books = await BooksAPI.getAll();
            const reader = AppState.allReaders.find(r => r.id === readerID);
            const currentlyReadingIds = reader.currently_reading.map(b => b.id);

            // Filter out books already being read
            const availableBooks = books.filter(b => !currentlyReadingIds.includes(b.id));

            if (availableBooks.length === 0) {
                UI.showNotification('All books are already in this reader\'s list!', 'info');
                return;
            }

            // Create modal HTML with search
            const modalHTML = `
                <div id="add-book-modal" class="modal">
                    <div class="modal-content modal-wide">
                        <div class="modal-header">
                            <h3>Add Book to Reading List</h3>
                            <span class="modal-close" onclick="Readers.closeAddBookModal()">&times;</span>
                        </div>
                        <div class="modal-search">
                            <input type="text" id="add-book-search" placeholder="Search books by title..." class="search-input">
                        </div>
                        <div class="books-selection" id="books-selection-list">
                            ${availableBooks.map(book => `
                                <div class="book-option" data-book-title="${escapeHtml(book.title).toLowerCase()}" onclick="Readers.addBook(${readerID}, ${book.id})">
                                    <div class="book-option-content">
                                        <strong>${escapeHtml(book.title)}</strong>
                                        <span class="book-description">${escapeHtml(book.description || 'No description')}</span>
                                        <span class="book-author">by ${escapeHtml(book.username)}</span>
                                    </div>
                                    <span class="book-add-icon">+</span>
                                </div>
                            `).join('')}
                        </div>
                    </div>
                </div>
            `;

            document.body.insertAdjacentHTML('beforeend', modalHTML);

            // Add search functionality
            const searchInput = document.getElementById('add-book-search');
            searchInput.addEventListener('input', (e) => {
                const searchTerm = e.target.value.toLowerCase();
                const bookOptions = document.querySelectorAll('.book-option');

                bookOptions.forEach(option => {
                    const title = option.getAttribute('data-book-title');
                    if (title.includes(searchTerm)) {
                        option.style.display = 'flex';
                    } else {
                        option.style.display = 'none';
                    }
                });
            });
        } catch (error) {
            UI.showNotification(error.message, 'error');
        }
    },

    closeAddBookModal() {
        const modal = document.getElementById('add-book-modal');
        if (modal) modal.remove();
    },

    async addBook(readerID, bookID) {
        try {
            await apiRequest(`/readers/${readerID}/books/${bookID}`, {
                method: 'POST'
            });

            this.closeAddBookModal();
            UI.showNotification('Book added to reading list!', 'success');
            await this.load();
        } catch (error) {
            UI.showNotification(error.message, 'error');
        }
    },

    async removeBook(readerID, bookID) {
        UI.showConfirmModal(
            'Remove this book from reading list?',
            async () => {
                try {
                    await apiRequest(`/readers/${readerID}/books/${bookID}`, {
                        method: 'DELETE'
                    });

                    UI.showNotification('Book removed from reading list!', 'success');
                    await this.load();
                } catch (error) {
                    UI.showNotification(error.message, 'error');
                }
            }
        );
    }
};
