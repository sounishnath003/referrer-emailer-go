const express = require('express');
const puppeteer = require('puppeteer');
const bodyParser = require('body-parser');

const app = express();
app.use(bodyParser.text({ limit: '5mb' }));

app.post('/generate-pdf', async (req, res) => {
    const html = req.body;

    console.log(`[${new Date().toISOString()}] Received request to /generate-pdf`);

    if (!html) {
        console.warn(`[${new Date().toISOString()}] Missing HTML content in request body`);
        return res.status(400).send('Missing HTML content');
    }

    try {
        console.log(`[${new Date().toISOString()}] Launching Puppeteer browser`);
        const browser = await puppeteer.launch({
            headless: 'new',
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });
        const page = await browser.newPage();
        console.log(`[${new Date().toISOString()}] Setting page content`);
        await page.setContent(html, { waitUntil: 'networkidle0' });

        console.log(`[${new Date().toISOString()}] Generating PDF`);
        const pdfBuffer = await page.pdf({
            format: 'Letter',
            printBackground: true,
            displayHeaderFooter: true,
            // margin: { top: '1in', bottom: '1in', left: '1in', right: '1in' },
            margin: { top: '0.12in', bottom: '0.12in', left: '0.15in', right: '0.15in' }
        });

        await browser.close();
        console.log(`[${new Date().toISOString()}] PDF generated and browser closed`);

        res.set({
            'Content-Type': 'application/pdf',
            'Content-Disposition': 'attachment; filename=resume.pdf',
        });

        res.send(pdfBuffer);
        console.log(`[${new Date().toISOString()}] PDF sent to client`);
    } catch (err) {
        console.error(`[${new Date().toISOString()}] Error generating PDF:`, err);
        res.status(500).send('Failed to generate PDF');
    }
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => console.log(`[${new Date().toISOString()}] PDF service running on port ${PORT}`));
