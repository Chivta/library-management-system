const Statistics = {
    async load() {
        try {
            const books = await BooksAPI.getAll();
            const readers = await ReadersAPI.getAll();

            const today = new Date().toDateString();
            const booksToday = books.filter(b => new Date(b.CreatedAt).toDateString() === today).length;
            const readersToday = readers.filter(r => new Date(r.CreatedAt).toDateString() === today).length;

            document.getElementById('stat-total-books').textContent = books.length;
            document.getElementById('stat-total-readers').textContent = readers.length;
            document.getElementById('stat-books-today').textContent = booksToday;
            document.getElementById('stat-readers-today').textContent = readersToday;
        } catch (error) {
            UI.showNotification('Failed to load statistics', 'error');
        }
    }
};
