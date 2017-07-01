import { Component, OnInit, ViewEncapsulation, Input } from '@angular/core';

@Component({
  selector: 'app-im-head',
  templateUrl: './im-head.component.html',
  styleUrls: ['./im-head.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHeadComponent implements OnInit {
  @Input() name;
  randomMatch: Array<HeadColorMatch> = [
    { bg: '#fff', text: '#000' },
    { bg: '#d05454', text: '#fff' },
    { bg: '#6dd067', text: '#fff' },
    { bg: '#676fd0', text: '#fff' },
    { bg: '#e47ae1', text: '#fff' },
    { bg: '#67c1d0', text: '#fff' },
    { bg: '#000', text: '#fff' },
    { bg: '#ffef2d', text: '#000' },
    { bg: '#eaae27', text: '#fff' },
    { bg: '#fbd1dc', text: '#000' },
  ]
  default: HeadColorMatch = { bg: '#fff', text: '#000' };
  constructor() { }

  ngOnInit() {
    if (this.name !== '') {
      this.name = this.name.substr(0, 1);
      this.default = this.randomMatch[this.getRandomArbitrary(0, 9)];
    }
  }
  getRandomArbitrary(min, max) {
    return Math.floor(Math.random() * (max - min) + min);
  }
}

export interface HeadColorMatch {
  bg?: string;
  text?: string;
}
