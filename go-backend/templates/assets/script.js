document.addEventListener('DOMContentLoaded', function () {
    var form = document.getElementById('urlForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();
        var originalUrl = document.getElementById('originalUrl').value;

        var linkElement = document.createElement('a');

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
                document.querySelector(".result-container").style.display = "block";
                addToHistory(shortenedUrl);
            })
            .catch(function (error) {
                console.error('Ошибка при отправке запроса:', error);
            });
    });
});

function addToHistory(url) {
    var historyContainer = document.querySelector(".history");
    var historyItem = document.createElement("div");
    historyItem.textContent = "Сокращенная ссылка: " + url;
    historyContainer.appendChild(historyItem);
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
    document.getElementById('shortenedUrlContainer').style.display = 'block';
});