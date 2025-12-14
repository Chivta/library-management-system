const Books = {
    applyFilters() {
        const search = AppState.bookFilters.search.toLowerCase();

        AppState.filteredBooks = AppState.allBooks.filter(book => {
            const matchesSearch = !search ||
                book.title.toLowerCase().includes(search) ||
                (book.description && book.description.toLowerCase().includes(search));
            return matchesSearch;
        });

        const [field, direction] = AppState.bookFilters.sortBy.split('-');
        AppState.filteredBooks.sort((a, b) => {
            let aVal, bVal;

            if (field === 'id') {
                aVal = a.id;
                bVal = b.id;
            } else if (field === 'title') {
                aVal = a.title.toLowerCase();
                bVal = b.title.toLowerCase();
            }

            if (direction === 'asc') {
                return aVal > bVal ? 1 : -1;
            } else {
                return aVal < bVal ? 1 : -1;
            }
        });

        AppState.booksCurrentPage = 1;
        this.render();
    },

    render() {
        const container = document.getElementById('books-list');

        if (!AppState.filteredBooks || AppState.filteredBooks.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3>No Books Found</h3>
                    <p>Start by adding your first book or adjust your filters!</p>
                </div>
            `;
            document.getElementById('books-pagination').innerHTML = '';
            return;
        }

        // Pagination logic
        const startIndex = (AppState.booksCurrentPage - 1) * AppState.booksItemsPerPage;
        const endIndex = startIndex + AppState.booksItemsPerPage;
        const paginatedBooks = AppState.filteredBooks.slice(startIndex, endIndex);

        const html = `
            <div class="items-grid">
                ${paginatedBooks.map(book => {
                    const canEdit = AppState.currentUser &&
                                   (book.user_id === AppState.currentUser.id || AppState.currentUser.role === 'admin');
                    return `
                        <div class="item-card" data-book-id="${book.id}">
                            <div class="item-card-header">
                                <span class="item-id">#${book.id}</span>
                                <span class="item-owner">by ${escapeHtml(book.username || 'Unknown')}</span>
                            </div>
                            <h4>${escapeHtml(book.title)}</h4>
                            <p>${escapeHtml(book.description || 'No description')}</p>
                            <div class="item-card-actions">
                                ${canEdit ? `
                                    <button class="btn btn-primary btn-small" onclick="Books.edit(${book.id})">Edit</button>
                                    <button class="btn btn-danger btn-small" onclick="Books.confirmDelete(${book.id})">Delete</button>
                                ` : ''}
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `;

        container.innerHTML = html;
        this.renderPagination();
    },

    renderPagination() {
        const container = document.getElementById('books-pagination');
        const totalPages = Math.ceil(AppState.filteredBooks.length / AppState.booksItemsPerPage);

        if (totalPages <= 1) {
            container.innerHTML = '';
            return;
        }

        let html = '<div class="pagination-controls">';
        html += `<button class="btn btn-small" onclick="Books.changePage(${AppState.booksCurrentPage - 1})" ${AppState.booksCurrentPage === 1 ? 'disabled' : ''}>← Previous</button>`;
        html += '<div class="pagination-pages">';

        for (let i = 1; i <= totalPages; i++) {
            if (i === 1 || i === totalPages || (i >= AppState.booksCurrentPage - 2 && i <= AppState.booksCurrentPage + 2)) {
                html += `<button class="btn btn-small ${i === AppState.booksCurrentPage ? 'btn-primary' : ''}" onclick="Books.changePage(${i})">${i}</button>`;
            } else if (i === AppState.booksCurrentPage - 3 || i === AppState.booksCurrentPage + 3) {
                html += '<span>...</span>';
            }
        }

        html += '</div>';
        html += `<button class="btn btn-small" onclick="Books.changePage(${AppState.booksCurrentPage + 1})" ${AppState.booksCurrentPage === totalPages ? 'disabled' : ''}>Next →</button>`;
        html += '</div>';

        container.innerHTML = html;
    },

    changePage(page) {
        const totalPages = Math.ceil(AppState.filteredBooks.length / AppState.booksItemsPerPage);
        if (page >= 1 && page <= totalPages) {
            AppState.booksCurrentPage = page;
            this.render();
            document.getElementById('books-list').scrollIntoView({ behavior: 'smooth' });
        }
    },

    async load() {
        const loading = document.getElementById('books-loading');
        UI.showLoading(loading);

        try {
            AppState.allBooks = await BooksAPI.getAll();
            AppState.filteredBooks = [...AppState.allBooks];
            this.applyFilters();
        } catch (error) {
            UI.showNotification(error.message, 'error');
            AppState.allBooks = [];
            AppState.filteredBooks = [];
            this.render();
        } finally {
            UI.hideLoading(loading);
        }
    },

    showModal(isEdit = false) {
        const modal = document.getElementById('book-modal');
        const title = document.getElementById('book-modal-title');

        title.textContent = isEdit ? 'Edit Book' : 'Add New Book';
        AppState.currentFormStep = 1;
        this.updateFormStep();
        UI.clearFormErrors('book');

        if (!isEdit) {
            document.getElementById('book-form').reset();
            AppState.editingBookId = null;
        }

        modal.classList.remove('hidden');
    },

    hideModal() {
        document.getElementById('book-modal').classList.add('hidden');
        AppState.editingBookId = null;
    },

    updateFormStep() {
        document.querySelectorAll('.step').forEach(step => {
            const stepNum = parseInt(step.dataset.step);
            if (stepNum === AppState.currentFormStep) {
                step.classList.add('active');
            } else if (stepNum < AppState.currentFormStep) {
                step.classList.add('completed');
                step.classList.remove('active');
            } else {
                step.classList.remove('active', 'completed');
            }
        });

        document.querySelectorAll('.form-step').forEach(step => {
            const stepNum = parseInt(step.dataset.step);
            if (stepNum === AppState.currentFormStep) {
                step.classList.add('active');
            } else {
                step.classList.remove('active');
            }
        });

        const prevBtn = document.getElementById('book-form-prev');
        const nextBtn = document.getElementById('book-form-next');
        const submitBtn = document.getElementById('book-form-submit');

        if (AppState.currentFormStep === 1) {
            prevBtn.classList.add('hidden');
        } else {
            prevBtn.classList.remove('hidden');
        }

        if (AppState.currentFormStep === 3) {
            nextBtn.classList.add('hidden');
            submitBtn.classList.remove('hidden');
            this.updateReviewStep();
        } else {
            nextBtn.classList.remove('hidden');
            submitBtn.classList.add('hidden');
        }
    },

    updateReviewStep() {
        document.getElementById('review-title').textContent = document.getElementById('book-title').value || '-';
        document.getElementById('review-category').textContent = document.getElementById('book-category').value || '-';
        document.getElementById('review-description').textContent = document.getElementById('book-description').value || '-';
        document.getElementById('review-date').textContent = document.getElementById('book-published-date').value || '-';
    },

    nextFormStep() {
        if (AppState.currentFormStep < 3) {
            AppState.currentFormStep++;
            this.updateFormStep();
        }
    },

    prevFormStep() {
        if (AppState.currentFormStep > 1) {
            AppState.currentFormStep--;
            this.updateFormStep();
        }
    },

    async edit(id) {
        try {
            const book = await BooksAPI.getById(id);
            this.showModal(true);

            document.getElementById('book-id').value = book.id;
            document.getElementById('book-title').value = book.title;
            document.getElementById('book-description').value = book.description || '';

            AppState.editingBookId = id;
        } catch (error) {
            UI.showNotification(error.message, 'error');
        }
    },

    async handleFormSubmit(e) {
        e.preventDefault();

        const submitBtn = document.getElementById('book-form-submit');
        submitBtn.disabled = true;
        submitBtn.textContent = 'Saving...';

        const bookData = {
            title: document.getElementById('book-title').value.trim(),
            description: document.getElementById('book-description').value.trim()
        };

        try {
            let message = '';
            if (AppState.editingBookId) {
                await BooksAPI.update(AppState.editingBookId, bookData);
                message = 'Book updated successfully';
            } else {
                await BooksAPI.create(bookData);
                message = 'Book created successfully';
            }

            this.hideModal();
            await this.load();
            UI.showNotification(message, 'success');
        } catch (error) {
            if (error.validationErrors) {
                UI.displayValidationErrors(error.validationErrors, 'book');
                UI.showNotification('Please fix the validation errors', 'error');
            } else {
                UI.showNotification(error.message, 'error');
            }
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Save Book';
        }
    },

    confirmDelete(id) {
        UI.showConfirmModal(
            'Are you sure you want to delete this book? This action cannot be undone.',
            async () => {
                try {
                    await BooksAPI.delete(id);
                    UI.showNotification('Book deleted successfully', 'success');
                    await this.load();
                } catch (error) {
                    UI.showNotification(error.message, 'error');
                }
            }
        );
    },

    confirmDeleteAll() {
        UI.showConfirmModal(
            'Are you sure you want to delete ALL books? This action cannot be undone and will remove all books from the database.',
            async () => {
                try {
                    await BooksAPI.deleteAll();
                    UI.showNotification('All books deleted successfully', 'success');
                    await this.load();
                } catch (error) {
                    UI.showNotification(error.message, 'error');
                }
            }
        );
    },

    export() {
        const csv = convertBooksToCSV(AppState.filteredBooks);
        downloadCSV(csv, 'books.csv');
        UI.showNotification('Books exported successfully', 'success');
    }
};
