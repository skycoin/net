import { Directive, ElementRef, Input, HostListener, Output, EventEmitter } from '@angular/core';

declare var Clipboard: any;
// tslint:disable-next-line:directive-selector
@Directive({ selector: '[clipboard]' })
export class ClipboardDirective {
  @Input() clipboardText = '';
  @Output() copySuccess = new EventEmitter<boolean>();
  cliboard: any = null;
  constructor(private el: ElementRef) {
    this.cliboard = new Clipboard(this.el.nativeElement);
    this.cliboard.on('success', (e) => {
      this.copySuccess.emit(true);
      e.clearSelection();
    });
  }

  @HostListener('click', ['$event'])
  _click(ev: Event) {
    this.cliboard.text = () => {
      return this.clipboardText;
    };
  }
}
