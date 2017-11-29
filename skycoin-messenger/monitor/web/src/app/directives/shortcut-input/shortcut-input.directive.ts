import { Directive, HostBinding, Input, ElementRef, OnInit, HostListener, Output, EventEmitter } from '@angular/core';

// tslint:disable-next-line:directive-selector
@Directive({ selector: '[shortcut-input]' })
export class ShortcutInputDirective implements OnInit {
  @Input() text = '';
  @HostBinding('class') classes = 'shortcut_input';
  @Output() onEdit = new EventEmitter<number>();
  isEdit = false;
  constructor(private el: ElementRef) { }
  ngOnInit() {
    this.el.nativeElement.value = this.text;
  }
  @HostListener('focus', ['$event'])
  _foucs(ev: Event) {
    this.isEdit = true;
  }
  @HostListener('blur', ['$event'])
  _blur(ev: Event) {
    if (this.isEdit) {
      this.isEdit = false;
      this.onEdit.emit(this.el.nativeElement.value);
    }
  }
}
