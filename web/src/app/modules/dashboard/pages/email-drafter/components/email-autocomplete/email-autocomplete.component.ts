import { NgIf } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-email-autocomplete',
  imports: [NgIf],
  templateUrl: './email-autocomplete.component.html',
  styleUrl: './email-autocomplete.component.css'
})
export class EmailAutocompleteComponent {
  @Input() suggestions: { email: string, companyName: string }[] = [];
  @Output() suggestionSelected = new EventEmitter<{ email: string, companyName: string }>();

  selectSuggestion(suggestion: { email: string, companyName: string }) {
    this.suggestionSelected.emit(suggestion);
  }
}
