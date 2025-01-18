import { NgFor, NgForOf } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  imports: [RouterLink, NgFor, NgForOf],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  templatesInformations = [
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
      shortDesc: "Craft an email asking job oppotunities at companies."
    },
    {
      label: "Send Congratulations",
      shortDesc: "Craft send congratulations roles customized email."
    },
    {
      label: "Appreciate Author",
      shortDesc: "Send appreciation note to Sounish Nath for this easy emailer service."
    },
  ]
}
