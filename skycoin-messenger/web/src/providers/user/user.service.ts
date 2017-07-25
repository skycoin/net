import { Injectable } from '@angular/core';
import { HeadColorMatch } from '../socket/msg'

@Injectable()
export class UserService {
  randomMatch: Array<HeadColorMatch> = [
    { bg: '#fff', color: '#000' },
    { bg: '#d05454', color: '#fff' },
    { bg: '#6dd067', color: '#fff' },
    { bg: '#676fd0', color: '#fff' },
    { bg: '#e47ae1', color: '#fff' },
    { bg: '#67c1d0', color: '#fff' },
    { bg: '#000', color: '#fff' },
    { bg: '#ffef2d', color: '#000' },
    { bg: '#eaae27', color: '#fff' },
    { bg: '#fbd1dc', color: '#000' },
  ]
  constructor() { }

  getRandomMatch() {
    return this.randomMatch[Math.floor(Math.random() * (9 - 0) + 0)];
  }
}
