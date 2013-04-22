package main

var Templates = `
{{define "objc/properties_m"}}
{{range . }}{{$f := .ToLanguageField "objc"}}@synthesize {{$f.Name | title}};
{{end}}
- (id) initWithDictionary:(NSDictionary*)dict{
	self = [super init];
	if (!self) {
		return self;
	}
	if (![dict isKindOfClass:[NSDictionary class]]) {
		return self;
	}
{{range . }}{{$f := .ToLanguageField "objc"}}{{ $name := $f.Name | title }}{{if .IsError}}	[self set{{$name}}:[{{$f.PkgName | title}} errorWithDictionary:{{$f.SetPropertyObjc}}]];{{else}}{{if $f.Primitive }}{{if $f.IsArray}}	[self set{{$name}}:{{$f.SetPropertyObjc}}];{{else}}	[self set{{$name}}:{{$name | $f.SetPropertyFromObjcDict}}];{{end}}{{else}}{{if $f.IsArray}}
	NSMutableArray * m{{$name}} = [[NSMutableArray alloc] init];
	NSArray * l{{$name}} = [dict valueForKey:@"{{$name}}"];
	if ([l{{$name}} isKindOfClass:[NSArray class]]) {
		for (NSDictionary * d in l{{$name}}) {
			[m{{$name}} addObject: [[{{$f.ConstructorType}} alloc] initWithDictionary:d]];
		}
	}
	[self set{{$name}}:m{{$name}}];{{else}}	[self set{{$name}}:[[{{$f.ConstructorType}} alloc] initWithDictionary:{{$name | $f.SetPropertyFromObjcDict}}]];{{end}}{{end}}{{end}}
{{end}}
	return self;
}

- (NSDictionary*) dictionary {
	NSMutableDictionary * dict = [[NSMutableDictionary alloc] init];
{{range . }}{{$f := .ToLanguageField "objc"}}{{ $name := $f.Name | title }}{{if $f.Primitive }}{{if $f.IsArray}}	[dict setValue:{{$f.GetPropertyObjc}} forKey:@"{{$name}}"];{{else}}	[dict setValue:{{$f.GetPropertyObjc | $f.GetPropertyToObjcDict}} forKey:@"{{$name}}"];{{end}}{{else}}{{if $f.IsArray}}
	NSMutableArray * m{{$name}} = [[NSMutableArray alloc] init];
	for ({{$f.Type}} p in {{$name}}) {
		[m{{$name}} addObject:[p dictionary]];
	}
	[dict setValue:m{{$name}} forKey:@"{{$name}}"];{{else}}	[dict setValue:[self.{{$name}} dictionary] forKey:@"{{$name}}"];{{end}}
	{{end}}
{{end}}
	return dict;
}
{{end}}

{{define "objc/properties_h"}}
{{range .}}{{$f := .ToLanguageField "objc"}}@property {{$f.PropertyAnnotation}} {{$f.FullObjcTypeName}} {{.Name | title}};
{{end}}
- (id) initWithDictionary:(NSDictionary*)dict;
- (NSDictionary*) dictionary;
{{end}}

{{define "objc/h"}}// Generated by github.com/sunfmin/goapigen
// DO NOT EDIT

#import <Foundation/Foundation.h>
{{$pkgName := .Name | title}}
@interface {{$pkgName}} : NSObject
@property (nonatomic, strong) NSString * BaseURL;
@property (nonatomic, assign) BOOL Verbose;
+ ({{$pkgName}} *) get;
@end

@interface Validated : NSObject
@end

{{range .DataObjects}}{{$do := .}}
// --- {{.Name}} ---
@interface {{.Name}} : NSObject
{{template "objc/properties_h" $do.Fields}}
@end
{{end}}

// === Interfaces ===
{{range .Interfaces}}{{$interface := .}}
{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
// --- {{.Name}}Params ---
@interface {{$interface.Name}}{{.Name}}Params : NSObject
{{template "objc/properties_h" .Params}}
@end

// --- {{.Name}}Results ---
@interface {{$interface.Name}}{{.Name}}Results : NSObject
{{template "objc/properties_h" .Results}}
@end
{{end}}{{end}}

@interface {{.Name}} : NSObject{{with .Constructor}}
{{template "objc/properties_h" .Method.Params}}
{{else}}
- (NSDictionary*) dictionary;
{{end}}
{{range .Methods}}
- {{with .ConstructorForInterface}}({{.Name}} *){{else}}({{$interface.Name | .ResultsForObjcFunction}}){{end}} {{.ParamsForObjcFunction}};
{{end}}@end
{{end}}
{{end}}

{{define "objc/m"}}// Generated by github.com/sunfmin/goapigen
// DO NOT EDIT
{{$pkgName := .Name | title}}
#import "{{.Name}}.h"

static {{.Name | title}} * _{{.Name}};
static NSDateFormatter * _dateFormatter;

@implementation {{.Name | title}} : NSObject
+ ({{.Name | title}} *)get {
	if(!_{{.Name}}) {
		_{{.Name}} = [[{{.Name | title}} alloc] init];
	}
	return _{{.Name}};
}

+ (NSDateFormatter *)dateFormatter {
	if(!_dateFormatter) {
		_dateFormatter = [[NSDateFormatter alloc] init];
		[_dateFormatter setDateFormat:@"yyyy-MM-dd'T'HH:mm:ss.SSSZZZZZ"];
	}
	return _dateFormatter;
}

+ (NSDictionary *) request:(NSURL*)url req:(NSDictionary *)req error:(NSError **)error {
	NSMutableURLRequest *httpRequest = [NSMutableURLRequest requestWithURL:url];
	[httpRequest setHTTPMethod:@"POST"];
	[httpRequest setValue:@"application/json;charset=utf-8" forHTTPHeaderField:@"Content-Type"];
	{{$pkgName}} * _api = [{{$pkgName}} get];
	NSData *requestBody = [NSJSONSerialization dataWithJSONObject:req options:NSJSONWritingPrettyPrinted error:error];
	if([_api Verbose]) {
		NSLog(@"Request: %@", [NSString stringWithUTF8String:[requestBody bytes]]);
	}
	[httpRequest setHTTPBody:requestBody];
	if(*error != nil) {
		return nil;
	}
	NSURLResponse  *response = nil;
	NSData *returnData = [NSURLConnection sendSynchronousRequest:httpRequest returningResponse:&response error:error];
	if(*error != nil || returnData == nil) {
		return nil;
	}
	if([_api Verbose]) {
		NSLog(@"Response: %@", [NSString stringWithUTF8String:[returnData bytes]]);
	}
	return [NSJSONSerialization JSONObjectWithData:returnData options:NSJSONReadingAllowFragments error:error];
}

+ (NSError *)errorWithDictionary:(NSDictionary *)dict {
	if (![dict isKindOfClass:[NSDictionary class]]) {
		return nil;
	}
	if ([[dict allKeys] count] == 0) {
		return nil;
	}
	NSMutableDictionary *userInfo = [NSMutableDictionary alloc];
	id reason = [dict valueForKey:@"Reason"];
	if ([reason isKindOfClass:[NSDictionary class]]) {
		userInfo = [userInfo initWithDictionary:reason];
	} else {
		userInfo = [userInfo init];
	}
	[userInfo setObject:[dict valueForKey:@"Message"] forKey:NSLocalizedDescriptionKey];

	NSString *code = [dict valueForKey:@"Code"];
	NSNumberFormatter *f = [[NSNumberFormatter alloc] init];
	[f setNumberStyle:NSNumberFormatterDecimalStyle];
	NSNumber *codeNumber = [f numberFromString:code];
	NSInteger intCode = -1;
	if (codeNumber != nil) {
		intCode = [codeNumber integerValue];
	}
	NSError *err = [NSError errorWithDomain:@"{{$pkgName}}Error" code:intCode userInfo:userInfo];
	return err;
}

@end

{{range .DataObjects}}{{$do := .}}
// --- {{.Name}} ---
@implementation {{.Name}}
{{template "objc/properties_m" $do.Fields}}
@end
{{end}}

// === Interfaces ===

{{range .Interfaces}}{{$interface := .}}
{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
// --- {{.Name}}Params ---
@implementation {{$interface.Name}}{{.Name}}Params : NSObject
{{template "objc/properties_m" .Params}}
@end

// --- {{.Name}}Results ---
@implementation {{$interface.Name}}{{.Name}}Results : NSObject
{{template "objc/properties_m" .Results}}
@end
{{end}}{{end}}{{end}}

{{range .Interfaces}}{{$interface := .}}
@implementation {{.Name}} : NSObject
{{with .Constructor}}
{{template "objc/properties_m" .Method.Params}}
{{else}}
- (NSDictionary*) dictionary {
	return [NSDictionary dictionaryWithObjectsAndKeys:nil];
}
{{end}}
{{range .Methods}}{{$method := .}}
// --- {{.Name}} ---
- {{with .ConstructorForInterface}}({{.Name}} *){{else}}({{$interface.Name | .ResultsForObjcFunction}}){{end}} {{.ParamsForObjcFunction}} {
	{{with .ConstructorForInterface}}
	{{.Name}} *results = [{{.Name}} alloc];
	{{range .Constructor.Method.Params}}{{$f := .ToLanguageField "objc"}}[results set{{$f.Name | title}}:{{$f.Name}}];
	{{end}}{{else}}
	{{$interface.Name}}{{.Name}}Results *results = [{{$interface.Name}}{{.Name}}Results alloc];
	{{$interface.Name}}{{.Name}}Params *params = [[{{$interface.Name}}{{.Name}}Params alloc] init];
	{{range .Params}}{{$f := .ToLanguageField "objc"}}[params set{{$f.Name | title}}:{{$f.Name}}];
	{{end}}
	{{$pkgName}} * _api = [{{$pkgName}} get];
	NSURL *url = [NSURL URLWithString:[NSString stringWithFormat:@"%@/{{$interface.Name}}/{{.Name}}.json", [_api BaseURL]]];
	if([_api Verbose]) {
		NSLog(@"Requesting URL: %@", url);
	}
	NSError *error;
	NSDictionary * dict = [{{$pkgName}} request:url req:[NSDictionary dictionaryWithObjectsAndKeys: [self dictionary], @"This", [params dictionary], @"Params", nil] error:&error];
	if(error != nil) {
		if([_api Verbose]) {
			NSLog(@"Error: %@", error);
		}
		results = [results init];
		[results setErr:error];
		return {{$method.ObjcReturnResultsOrOnlyOne}};
	}
	results = [results initWithDictionary: dict];
	{{end}}
	return {{$method.ObjcReturnResultsOrOnlyOne}};
}
{{end}}@end
{{end}}


{{end}}


{{define "httpserver"}}// Generated by github.com/sunfmin/goapigen
// DO NOT EDIT
{{$pkg := .}}
package {{.Name}}httpimpl

import ({{range .ServerImports}}
	"{{.}}"{{end}}
)

var _ govalidations.Errors
var _ = time.Sunday

type CodeError interface {
	Code() string
}

type SerializableError struct {
	Code    string
	Message string
	Reason  error
}

func (s *SerializableError) Error() string {
	return s.Message
}

func NewError(err error) (r error) {
	se := &SerializableError{Message:err.Error()}
	ce, yes := err.(CodeError)
	if yes {
		se.Code = ce.Code()
	}
	se.Reason = err
	r = se
	return
}

func AddToMux(prefix string, mux *http.ServeMux) {
	{{range .Interfaces}}{{$interface := .}}{{range .Methods}}{{if .ConstructorForInterface}}{{else}}
	mux.HandleFunc(prefix+"/{{$interface.Name}}/{{.Name}}.json", {{$interface.Name}}_{{.Name}}){{end}}{{end}}{{end}}
	return
}
{{range .Interfaces}}{{$interface := .}}
{{with .Constructor}}{{else}}
var {{$interface.Name | downcase}} {{$pkg.Name}}.{{$interface.Name}} = {{$pkg.ImplPkg | dotlastname}}.Default{{$interface.Name}}{{end}}

type {{$interface.Name}}Data struct {
{{with .Constructor}}{{range .Method.Params}}	{{.Name | title}} {{.FullGoTypeName}}
{{end}}{{end}}}

{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
type {{$interface.Name}}_{{$method.Name}}_Params struct {
{{with $interface.Constructor}}	This   {{$interface.Name}}Data
{{end}}	Params struct {
{{range .Params}}		{{.Name | title}} {{.FullGoTypeName}}
{{end}}	}
}

type {{$interface.Name}}_{{$method.Name}}_Results struct {
{{range .Results}}	{{.Name | title}} {{.FullGoTypeName}}
{{end}}
}

func {{$interface.Name}}_{{$method.Name}}(w http.ResponseWriter, r *http.Request) {
	var p {{$interface.Name}}_{{$method.Name}}_Params
	if r.Body == nil {
		panic("no body")
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	var result {{$interface.Name}}_{{$method.Name}}_Results
	enc := json.NewEncoder(w)
	if err != nil {
		result.Err = NewError(err)
		enc.Encode(result)
		return
	}
{{if $interface.Constructor}}
	s, err := {{$interface.Constructor.FromInterface.Name | downcase }}.{{$interface.Constructor.Method.Name}}({{$interface.Constructor.Method.ParamsForGoServerConstructorFunction}})
{{else}}
	s := {{$interface.Name | downcase }}
{{end}}
	if err != nil {
		result.Err = NewError(err)
		enc.Encode(result)
		return
	}
	{{$method.ResultsForGoServerFunction "result"}} = s.{{$method.Name}}({{$method.ParamsForGoServerFunction}})
	if result.Err != nil {
		result.Err = NewError(result.Err)
	}
	err = enc.Encode(result)
	if err != nil {
		panic(err)
	}
	return
}
{{end}}{{end}}

{{end}}





{{end}}

{{define "javascript/interfaces"}}// Generated by github.com/sunfmin/goapigen
// DO NOT EDIT

(function(api, $, undefined ) {
	api.rpc = function(endpoint, input, callback) {
		var methodUrl = api.baseurl + endpoint;
		var message = JSON.stringify(input);
		var req = $.ajax({
			type: "POST",
			url: methodUrl,
			contentType:"application/json; charset=utf-8",
			dataType:"json",
			processData: false,
			data: message
		});
		req.done(function(data, textStatus, jqXHR) {
			callback(data);
		});
	};
})( window.{{.Name}} = window.{{.Name}} || {}, jQuery);



(function( api, undefined ) {
{{range .Interfaces}}{{ $interfaceName := .Name}}
	api.{{$interfaceName}} = function() {};
{{range .Methods}}{{$method := .}}{{if .ConstructorForInterface}}
	api.{{$interfaceName}}.prototype.{{.Name}} = function({{$method.ParamsForJavascriptFunction}}) {
		var r = new api.{{.ConstructorForInterface.Name}}(){{range .Params}};
		r.{{.Name | title}} = {{.Name}}{{end}};
		return r;
	}
{{else}}
	api.{{$interfaceName}}.prototype.{{.Name}} = function({{$method.ParamsForJavascriptFunction}}{{if $method.ParamsForJavascriptFunction}}, {{end}}callback) {
		api.rpc("/{{$interfaceName}}/{{.Name}}.json", {"This": this, "Params": {{$method.ParamsForJson}}}, function(data){
			callback({{$method.ResultsForJavascriptFunction "data"}})
		});
		return;
	}
{{end}}{{end}}{{end}}

}( window.{{.Name}} = window.{{.Name}} || {} ));

{{end}}

`
