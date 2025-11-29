let nextEntry = null;
let currentCalled = null;

// Обновление UI
function updateLabels() {
    const nextLabel = document.getElementById('next-label');
    const calledLabel = document.getElementById('called-label');

    if (nextEntry) {
        nextLabel.textContent = `Следующий, талон номер #${nextEntry.id}`;
    } else {
        nextLabel.textContent = 'Следующий: —';
    }

    if (currentCalled) {
        calledLabel.textContent = `Текущий вызванный, талон номер #${currentCalled.id}`;
    } else {
        calledLabel.textContent = 'Текущий вызванный: —';
    }
}

// Показать ошибку
function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.classList.add('show');
    setTimeout(() => {
        errorDiv.classList.remove('show');
    }, 5000);
}

// Получить следующего
async function refreshNext() {
    try {
        const response = await fetch('/api/next');
        const data = await response.json();

        if (data.available && data.entry) {
            nextEntry = data.entry;
        } else {
            nextEntry = null;
        }
        updateLabels();
    } catch (error) {
        showError('Ошибка при получении следующего: ' + error.message);
    }
}

// Позвать следующего
async function callNext() {
    try {
        const response = await fetch('/api/call-next', { method: 'POST' });
        const data = await response.json();

        if (data.empty) {
            alert('Очередь пуста.');
        } else if (data.entry) {
            currentCalled = data.entry;
            nextEntry = null;
        }
        updateLabels();
    } catch (error) {
        showError('Ошибка при вызове: ' + error.message);
    }
}

// Отметить как обслуженного
async function markDone() {
    if (!currentCalled) {
        alert('Сначала позовите кого-то.');
        return;
    }

    try {
        const response = await fetch(`/api/done/${currentCalled.id}`, { method: 'POST' });
        const data = await response.json();

        if (data.success) {
            alert(`Билет ${currentCalled.id} отмечен как 'done'.`);
            currentCalled = null;
            updateLabels();
        }
    } catch (error) {
        showError('Ошибка при отметке: ' + error.message);
    }
}

// Сбросить текущего
function clearCurrent() {
    currentCalled = null;
    updateLabels();
}

// События кнопок
document.getElementById('btn-refresh').addEventListener('click', refreshNext);
document.getElementById('btn-call').addEventListener('click', callNext);
document.getElementById('btn-done').addEventListener('click', markDone);
document.getElementById('btn-clear').addEventListener('click', clearCurrent);

// Автоматическое обновление каждые 3 секунды
setInterval(refreshNext, 3000);

// Первоначальная загрузка
refreshNext();
