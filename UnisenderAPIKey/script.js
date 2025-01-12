define(['jquery', 'underscore', 'twigjs'], function ($, _, twigjs) {
  var CustomWidget = function () {
    var self = this;
    var tokenURL = 'https://d555-172-241-70-196.ngrok-free.app';

    this.getTemplate = _.bind(function (template, params, callback) {
      params = typeof params == 'object' ? params : {};
      template = template || '';

      return this.render({
        href: '/templates/' + template + '.twig',
        base_path: this.params.path,
        v: this.get_version(),
        load: callback
      }, params);
    }, this);

    this.callbacks = {
      render: function () {
        return true;
      },

      init: _.bind(function () {
        console.log('init');

        AMOCRM.addNotificationCallback(self.get_settings().widget_code, function (data) {
          console.log(data);
        });

        this.add_action('phone', function (params) {
          /**
           * код взаимодействия с виджетом телефонии
           */
          console.log(params);
        });

        return true;
      }, this),

      bind_actions: function () {
        console.log('bind_actions');

        return true;
      },

      settings: _.bind(function ($modal_body) {
        this.getTemplate(
          'oferta',
          {},
          function (template) {
            return true;
          }
        );

        return true;
      }, this),

      onSave: function () {
      const input= 'input[name="api_token"]';
      var ttoken = jQuery(input).val();
      const obj ={unisender_key:ttoken,
       account_id:AMOCRM.constant('account').id};
        jQuery.ajax({
          method: 'POST',
          url: tokenURL,
          dataType: 'json',
          data: obj,
          always: function (data) {
            console.warn(data)
          }
        });

        return true;
      },

      destroy: function () {
        _.noop();
      },

      advancedSettings: _.bind(function () {
        var $work_area = $('#work-area-' + self.get_settings().widget_code),
            $save_button = $(
              twigjs({ ref: '/tmpl/controls/button.twig' }).render({
                text: 'Сохранить',
                class_name: 'button-input_blue button-input-disabled js-button-save-' + self.get_settings().widget_code,
                additional_data: ''
              })
            ),
            $cancel_button = $(
              twigjs({ ref: '/tmpl/controls/cancel_button.twig' }).render({
                text: 'Отмена',
                class_name: 'button-input-disabled js-button-cancel-' + self.get_settings().widget_code,
                additional_data: ''
              })
            );

        console.log('advancedSettings');

        $('.content__top__preset').css({ float: 'left' });

        $('.list__body-right__top').css({ display: 'block' })
          .append('<div class="list__body-right__top__buttons"></div>');
        $('.list__body-right__top__buttons').css({ float: 'right' })
          .append($cancel_button)
          .append($save_button);

        self.getTemplate('advanced_settings', {}, function (template) {
          var $page = $(
            template.render({ title: self.i18n('advanced').title, widget_code: self.get_settings().widget_code })
          );

          $work_area.append($page);
        });
      }, self),
    };

    return this;
  };

  return CustomWidget;
});
