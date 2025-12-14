function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function convertBooksToCSV(books) {
    const headers = ['ID', 'Title', 'Description'];
    const rows = books.map(book => [
        book.id,
        `"${book.title.replace(/"/g, '""')}"`,
        `"${(book.description || '').replace(/"/g, '""')}"`
    ]);

    return [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
}

function convertReadersToCSV(readers) {
    const headers = ['ID', 'Name', 'Surname'];
    const rows = readers.map(reader => [
        reader.id,
        `"${reader.name.replace(/"/g, '""')}"`,
        `"${reader.surname.replace(/"/g, '""')}"`
    ]);

    return [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
}

function downloadCSV(csv, filename) {
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    window.URL.revokeObjectURL(url);
}
