package sync

import (
	"reflect"
	"sync"
	"testing"
)

func TestConcurrentMap_LoadOrStoreV1(t *testing.T) {
	type fields struct {
		values map[string]string
		lock   sync.RWMutex
	}
	type args struct {
		key      string
		newValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "loadV1",
			fields: fields{
				values: map[string]string{},
				lock:   sync.RWMutex{},
			},
			args: args{
				key:      "key1",
				newValue: "value1",
			},
			want:  "value1",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			concurrentMap := &ConcurrentMap[string, string]{
				values: tt.fields.values,
				lock:   tt.fields.lock,
			}
			got, got1 := concurrentMap.LoadOrStoreV1(tt.args.key, tt.args.newValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadOrStoreV1() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadOrStoreV1() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConcurrentMap_LoadOrStoreV2(t *testing.T) {
	type fields struct {
		values map[string]string
		lock   sync.RWMutex
	}
	type args struct {
		key      string
		newValue string
	}
	values := map[string]string{}
	lock := sync.RWMutex{}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "loadV2-1",
			fields: fields{
				values: values,
				lock:   lock,
			},
			args: args{
				key:      "key1",
				newValue: "value1",
			},
			want:  "value1",
			want1: false,
		},
		{
			name: "loadV2-2",
			fields: fields{
				values: values,
				lock:   lock,
			},
			args: args{
				key:      "key1",
				newValue: "value1",
			},
			want:  "value1",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			concurrentMap := &ConcurrentMap[string, string]{
				values: tt.fields.values,
				lock:   tt.fields.lock,
			}
			got, got1 := concurrentMap.LoadOrStoreV2(tt.args.key, tt.args.newValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadOrStoreV2() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadOrStoreV2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConcurrentMap_LoadOrStoreV3(t *testing.T) {
	type fields struct {
		values map[string]string
		lock   sync.RWMutex
	}
	type args struct {
		key      string
		newValue string
	}
	values := map[string]string{}
	lock := sync.RWMutex{}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "loadV2-1",
			fields: fields{
				values: values,
				lock:   lock,
			},
			args: args{
				key:      "key1",
				newValue: "value1",
			},
			want:  "value1",
			want1: false,
		},
		{
			name: "loadV2-2",
			fields: fields{
				values: values,
				lock:   lock,
			},
			args: args{
				key:      "key1",
				newValue: "value1",
			},
			want:  "value1",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			concurrentMap := &ConcurrentMap[string, string]{
				values: tt.fields.values,
				lock:   tt.fields.lock,
			}
			got, got1 := concurrentMap.LoadOrStoreV3(tt.args.key, tt.args.newValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadOrStoreV3() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadOrStoreV3() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
