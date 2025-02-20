import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-subhero',
  imports: [],
  templateUrl: './subhero.component.html',
  styleUrl: './subhero.component.css'
})
export class SubheroComponent {
  @Input() title!: string;
  @Input() subtitle!: string;
}
