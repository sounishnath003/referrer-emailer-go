import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ProfileService } from '../../services/profile.service';
import { MarkdownModule } from 'ngx-markdown';
import { NgIf, DatePipe, } from '@angular/common';
import { FormsModule } from '@angular/forms';
import jsPDF from 'jspdf';
import { marked } from 'marked';

@Component({
    selector: 'app-tailored-resume-view',
    standalone: true,
    imports: [NgIf, MarkdownModule, DatePipe, FormsModule],
    templateUrl: './tailored-resume-view.component.html'
})
export class TailoredResumeViewComponent implements OnInit {
    resume: any = null;
    loading = true;
    error: string | null = null;
    editMode = false;
    editableMarkdown = '';

    constructor(private route: ActivatedRoute, private profileService: ProfileService) { }

    ngOnInit() {
        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.profileService.getTailoredResumeById$(id).subscribe({
                next: (res) => {
                    this.resume = res;
                    this.editableMarkdown = res.resumeMarkdown;
                    this.loading = false;
                },
                error: (err) => {
                    this.error = err.error?.error || 'Failed to load tailored resume.';
                    this.loading = false;
                }
            });
        } else {
            this.error = 'No resume ID provided.';
            this.loading = false;
        }
    }

    toggleEdit() {
        this.editMode = !this.editMode;
        // When toggling from edit mode to preview mode, update the original resume content and save to DB
        if (!this.editMode && this.resume) {
            this.resume.resumeMarkdown = this.editableMarkdown;
            this.loading = true;
            this.profileService.updateTailoredResumeMarkdown$(this.resume.id, this.editableMarkdown).subscribe({
                next: () => {
                    this.loading = false;
                },
                error: (err) => {
                    this.error = err.error?.error || 'Failed to save changes.';
                    this.loading = false;
                }
            });
        }
    }

    /**
     * Download the resume as a PDF, preserving the on-screen Tailwind/CSS styling.
     * This method captures the actual rendered HTML (with all styles applied) and converts it to PDF.
     * It uses html2pdf.js, which leverages html2canvas and jsPDF under the hood.
     * 
     * Note: html2pdf.js must be installed and available (npm install html2pdf.js).
     * If using Angular, you may need to import it dynamically to avoid SSR/build issues.
     */
    async downloadPDF() {
        // Dynamically import html2pdf.js to avoid SSR/build issues
        const html2pdf = (await import('html2pdf.js')).default;

        // Get the actual rendered resume content (with Tailwind/CSS applied)
        const element = document.getElementById('resume-pdf-content');
        if (!element) {
            this.error = 'Resume content not found for PDF export.';
            return;
        }

        // Optional: temporarily remove box-shadow or adjust styles for PDF clarity
        element.classList.remove('p-8');
        element.classList.add('print-pdf');
        element.style.fontSize = '12px';

        // Configure PDF options for resume look (Letter, margins, scale, etc)
        const opt = {
            margin: 0,
            filename: 'tailored-resume.pdf',
            html2canvas: { scale: 4, useCORS: true, backgroundColor: '#fff' },
            jsPDF: { unit: 'in', format: 'letter', orientation: 'portrait' },
            pagebreak: { mode: ['avoid-all', 'css', 'legacy'] }
        };

        // Use html2pdf to generate and save the PDF
        html2pdf().set(opt).from(element).save().finally(() => {
            element.classList.remove('print-pdf');
        });
    }
}
