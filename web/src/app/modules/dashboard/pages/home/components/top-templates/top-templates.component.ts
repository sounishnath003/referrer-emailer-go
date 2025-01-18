import { Component, Input } from '@angular/core';
import { TemplateInformationType } from '../types';

@Component({
  selector: 'app-top-templates',
  imports: [],
  templateUrl: './top-templates.component.html',
  styleUrl: './top-templates.component.css'
})
export class TopTemplateComponent {
  @Input() templatesInformation!: TemplateInformationType;
}
