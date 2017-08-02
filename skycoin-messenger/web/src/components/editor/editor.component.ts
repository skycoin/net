import { Component, OnInit, ViewEncapsulation, ViewChild, ElementRef, Output, EventEmitter } from '@angular/core';
import { Http } from '@angular/http';
import 'rxjs/add/operator/map'
import * as emojione from 'emojione';

declare const Buffer

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.components.scss'],
  encapsulation: ViewEncapsulation.None
})
export class EditorComponent implements OnInit {
  emojis = [
    { code: '1f600', short: ':grinning:', desc: 'Grinning face' },
    { code: '1f601', short: ':grin:', desc: 'Grinning face with smiling eyes' },
    { code: '1f602', short: ':joy:', desc: 'Face with tears of joy' },
    { code: '1f603', short: ':smiley:', desc: 'Smiling face with open mouth' },
    { code: '1f604', short: ':smile:', desc: 'Smiling face with open mouth and smiling eyes' },
    { code: '1f605', short: ':sweat_smile:', desc: 'Smiling face with open mouth and cold sweat' },
    { code: '1f606', short: ':laughing:', desc: 'Smiling face with open mouth and tightly-closed eyes' },
    { code: '1f607', short: ':innocent:', desc: 'Smiling face with halo' },

    { code: '1f608', short: ':smiling_imp:', desc: 'Smiling face with horns' },
    { code: '1f609', short: ':wink:', desc: 'Winking face' },
    { code: '1f60a', short: ':blush:', desc: 'Smiling face with smiling eyes' },
    { code: '1f60b', short: ':yum:', desc: 'Face savoring delicious food' },
    { code: '1f60c', short: ':relieved:', desc: 'Relieved face' },
    { code: '1f60d', short: ':heart_eyes:', desc: 'Smiling face with heart-shaped eyes' },
    { code: '1f60e', short: ':sunglasses:', desc: 'Smiling face with sunglasses' },
    { code: '1f60f', short: ':smirk:', desc: 'Smirking face' },

    { code: '1f610', short: ':neutral_face:', desc: 'Neutral face' },
    { code: '1f611', short: ':expressionless:', desc: 'Expressionless face' },
    { code: '1f612', short: ':unamused:', desc: 'Unamused face' },
    { code: '1f613', short: ':sweat:', desc: 'Face with cold sweat' },
    { code: '1f614', short: ':pensive:', desc: 'Pensive face' },
    { code: '1f615', short: ':confused:', desc: 'Confused face' },
    { code: '1f616', short: ':confounded:', desc: 'Confounded face' },
    { code: '1f617', short: ':kissing:', desc: 'Kissing face' },

    { code: '1f618', short: ':kissing_heart:', desc: 'Face throwing a kiss' },
    { code: '1f619', short: ':kissing_smiling_eyes:', desc: 'Kissing face with smiling eyes' },
    { code: '1f61a', short: ':kissing_closed_eyes:', desc: 'Kissing face with closed eyes' },
    { code: '1f61b', short: ':stuck_out_tongue:', desc: 'Face with stuck out tongue' },
    { code: '1f61c', short: ':stuck_out_tongue_winking_eye:', desc: 'Face with stuck out tongue and winking eye' },
    { code: '1f61d', short: ':stuck_out_tongue_closed_eyes:', desc: 'Face with stuck out tongue and tightly-closed eyes' },
    { code: '1f61e', short: ':disappointed:', desc: 'Disappointed face' },
    { code: '1f61f', short: ':worried:', desc: 'Worried face' },

    { code: '1f620', short: ':angry:', desc: 'Angry face' },
    { code: '1f621', short: ':rage:', desc: 'Pouting face' },
    { code: '1f622:', short: ':cry', desc: 'Crying face' },
    { code: '1f623', short: ':persevere:', desc: 'Persevering face' },
    { code: '1f624', short: ':triumph:', desc: 'Face with look of triumph' },
    { code: '1f625', short: ':disappointed_relieved:', desc: 'Disappointed but relieved face' },
    { code: '1f626', short: ':frowning:', desc: 'Frowning face with open mouth' },
    { code: '1f627', short: ':anguished:', desc: 'Anguished face' },

    { code: '1f628', short: ':fearful:', desc: 'Fearful face' },
    { code: '1f629', short: ':weary:', desc: 'Weary face' },
    { code: '1f62a', short: ':sleepy:', desc: 'Sleepy face' },
    { code: '1f62b', short: ':tired_face:', desc: 'Tired face' },
    { code: '1f62c', short: ':grimacing:', desc: 'Grimacing face' },
    { code: '1f62d', short: ':sob:', desc: 'Loudly crying face' },
    { code: '1f62e', short: ':open_mouth:', desc: 'Face with open mouth' },
    { code: '1f62f', short: ':hushed:', desc: 'Hushed face' },

    { code: '1f630', short: ':cold_sweat:', desc: 'Face with open mouth and cold sweat' },
    { code: '1f631', short: ':scream:', desc: 'Face screaming in fear' },
    { code: '1f632', short: ':astonished:', desc: 'Astonished face' },
    { code: '1f633', short: ':flushed:', desc: 'Flushed face' },
    { code: '1f634', short: ':sleeping:', desc: 'Sleeping face' },
    { code: '1f635', short: ':dizzy_face:', desc: 'Dizzy face' },
    { code: '1f636', short: ':no_mouth:', desc: 'Face without mouth' },
    { code: '1f637', short: ':mask:', desc: 'Face with medical mask' }];
  emojiBox = false;
  @Output() onEnter = new EventEmitter<string>();
  @ViewChild('editor') editor: ElementRef;
  constructor(private http: Http) { }

  ngOnInit() {
  }
  openEmoji(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.emojiBox = !this.emojiBox;
  }
  selectEmoji(ev: Event, code: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (code === '') {
      return;
    }
    this.edit(code);
  }

  _enter(e: Event) {
    e.preventDefault();
    e.stopImmediatePropagation();
    e.stopPropagation();
    let msg = '';
    const nodes = this.editor.nativeElement.childNodes
    nodes.forEach(el => {
      switch (el.nodeName) {
        case '#text':
          msg += el.data;
          break;
        case 'IMG':
          msg += el.title;
          break;
        default:
          break;
      }
    });
    if (Buffer.byteLength(msg, 'utf8') >= 500) {
      console.log('max length 500');
      return;
    }
    this.onEnter.emit(msg);
    this.editor.nativeElement.innerHTML = '';
  }
  close(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.emojiBox = false;
  }
  edit(code: string) {
    this.emojiBox = false;
    this.editor.nativeElement.focus();
    const text = this.editor.nativeElement.innerHTML;
    this.editor.nativeElement.innerHTML += emojione.shortnameToImage(code);
    const range = document.createRange();
    const sel = window.getSelection();
    range.selectNodeContents(this.editor.nativeElement);
    sel.removeAllRanges();
    sel.addRange(range);
    sel.collapseToEnd();
  }
}
