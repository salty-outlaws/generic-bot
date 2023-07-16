package plugin

import (
	"reflect"
	"testing"
)

func TestLoadPluginRepo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Test",
			args: args{
				url: "https://github.com/varunbheemaiah/generic-bot-plugins/blob/master/config.json?raw=true",
			},
			want: []string{
				"insult_compliment.lua",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPluginRepo(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPluginRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadPluginRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
