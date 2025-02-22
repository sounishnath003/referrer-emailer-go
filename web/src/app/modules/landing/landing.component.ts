import { UpperCasePipe } from '@angular/common';
import { Component } from '@angular/core';
import { FooterComponent } from "../../shared/components/footer/footer.component";

@Component({
  selector: 'app-landing',
  imports: [UpperCasePipe, FooterComponent],
  templateUrl: './landing.component.html',
  styleUrl: './landing.component.css'
})
export class LandingComponent {

}
