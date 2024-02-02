document.addEventListener('DOMContentLoaded', function () {
    const historyListContainer = document.getElementById('historyList');
    fetch('/get')
        .then(response => response.json())
        .then(data => {
            if (data.response && Object.keys(data.response).length > 0) {
                Object.entries(data.response).forEach(([fullLink, shortLink]) => {
                    const listItem = document.createElement('li');
                    listItem.textContent = `${fullLink} => ${shortLink}`;
                    historyListContainer.appendChild(listItem);
                });
            } else {
                const listItem = document.createElement('li');
                listItem.textContent = 'История сокращений пуста.';
                historyListContainer.appendChild(listItem);
            }
        })
        .catch(error => console.error('Ошибка при запросе данных:', error));

    var form = document.getElementById('urlForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();
        var originalUrl = document.getElementById('originalUrl').value;

        fetch('/shorter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: 'full_link=' + encodeURIComponent(originalUrl),
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (data) {
                var shortenedUrlContainer = document.getElementById('shortenedUrl');
                var shortenedUrl = data.response.message;
                shortenedUrlContainer.innerHTML = '';
                shortenedUrlContainer.textContent = shortenedUrl;
                shortenedUrlContainer.setAttribute('contentEditable', 'true');
                document.querySelector(".result-container").style.display = "flex";
                document.getElementById('shortenedUrlContainer').style.display = 'block';
            })
            .catch(function (error) {
                console.error('Ошибка при отправке запроса:', error);
            });
    });
});

function addToHistory(url) {
    var historyContainer = document.querySelector(".history");
    var listItem = document.createElement("div");
    listItem.textContent = url;
    historyContainer.appendChild(listItem);
}

function copyShortenedUrl() {
    var shortenedUrlField = document.getElementById("shortenedUrl");
    var range = document.createRange();
    range.selectNodeContents(shortenedUrlField);
    var selection = window.getSelection();
    selection.removeAllRanges();
    selection.addRange(range);
    document.execCommand("copy");
    selection.removeAllRanges();
    alert("Скопирована сокращенная ссылка: " + shortenedUrlField.textContent);
}

document.getElementById('urlForm').addEventListener('submit', function (event) {
    event.preventDefault();
    document.getElementById('shortenedUrlContainer').style.display = 'flex';
});
