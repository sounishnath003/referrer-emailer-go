import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { TopTemplateComponent } from "./components/top-templates/top-templates.component";
import { TemplateInformationType } from './components/types';
import { SubheroComponent } from '../shared/subhero/subhero.component';


@Component({
  selector: 'app-home',
  imports: [RouterLink, TopTemplateComponent, SubheroComponent],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  templatesInformations: TemplateInformationType[] = [
    {
      label: "Software Engineering",
      shortDesc: "Craft Software engineering roles customized email."
    },
    {
      label: "Data Engineering",
      shortDesc: "Craft Data engineering roles customized email."
    },
    {
      label: "Business Analyst",
      shortDesc: "Craft Business engineering roles customized email."
    },
    {
      label: "Ask for Job Opportunities",
      shortDesc: "Craft an email asking job opportunities at companies."
    },
    {
      label: "Send Congratulations",
      shortDesc: "Craft send congratulations roles customized email."
    },
    {
      label: "Draft with AI",
      shortDesc: "Use Generative AI to draft an customized message"
    },
  ]
}
