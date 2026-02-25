
const listContainer = document.getElementById('article-list');
const mainContainer = document.getElementById('content-area');
const dateContainer = document.getElementById('added-on-date');

function renderSidebar() {
    listContainer.innerHTML = '';
    myArticles.forEach((article, index) => {
        const li = document.createElement('li');
        const a = document.createElement('a');
        a.href = "#";
        a.textContent = article.title;
        a.onclick = (e) => {
            e.preventDefault();
            displayArticle(index);
        };

        li.appendChild(a);
        listContainer.appendChild(li);
    });
}
function displayArticle(index) {
    const article = myArticles[index];

    if (article) {
        // Update HTML Content
        mainContainer.innerHTML = article.content || "<em>No content available.</em>";

        // Update Date (The fix!)
        if (dateContainer) {
            const dateStr = article.metadata ? article.metadata.LastModifiedStr : "Unknown Date";
            dateContainer.textContent = `Added on: ${dateStr}`;
        }
    }
}

// 1. Build the list
renderSidebar();

// 2. Load the first article (Content + Date) immediately
if (myArticles.length > 0) {
    displayArticle(0);
}
