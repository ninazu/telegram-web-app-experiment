async function loadMenuData() {
    try {
        const response = await fetch('https://example.com/api/menu'); // Заміна URL на реальний
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        const data = await response.json();
        const htmlMenu = generateMenu(data.menu);
        document.getElementById('menu-container').innerHTML = htmlMenu;
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}