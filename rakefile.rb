# This file is use to generate Azure SDK from swagger. 
# Swaggers specs are located in git repo: https://github.com/Azure/azure-rest-api-specs
# AutoRest is located in git repo: https://github.com/Azure/autorest
# Build AutoRest before regenerating the SDK...
# cd local/AutoRest/git/repo
# gulp build
# Command to run rake file: 
# 				rake SDK_VERSION=<major>.<minor>.<patch> SWAGGER=path AUTOREST=path
# 				where SDK_VERSION is SDK version,
#               SWAGGER is local swagger specs git repo (defaults to the github repo),
#               AUTOREST is local AutoRest git repo 


require 'fileutils'

AUTOREST     = "%s/autorest/src/core/AutoRest/bin/Debug/net451/win7-x64/AutoRest.exe"

ARM_SWAGGERS = {
	# asazure: {version: "2016-05-06"},
	authorization: {version: "2015-07-01"},
	batch: {version: "2015-12-01", swagger: "BatchManagement"},
	cdn: {version: "2016-04-02"},
	cognitiveservices: {version: "2016-02-01-preview"},
	compute: {version: "2016-03-30"}, #composite swagger (includes container service)
	# containerservice: {version: "2016-03-30", swagger: "containerservice"},
	# commerce: {version: "2015-06-01-preview"},
	datalake_analytics: {
		account: {version: "2015-10-01-preview"},
		# catalog: {version: "2016-06-01-preview"},
		# job:  {version: "2016-03-20-preview"}
	},
	datalake_store: {
		account: {version: "2015-10-01-preview"},
		filesystem: {version: "2015-10-01-preview"}
	},
	devtestlabs: {version: "2016-05-15", swagger: "DTL"},
	dns: {version: "2016-04-01"},
	eventhub: {version: "2015-08-01", swagger: "EventHub"},
	# graphrbac: {version: "1.6"}, # composite swagger
	# insights (composite swagger)
	intune: {version: "2015-01-14-preview"},
	iothub: {version: "2016-02-03"},
	keyvault: {version: "2015-06-01"},
	logic: {version: "2016-06-01"}, #composite swagger
	# machine learning has two swaggers, but not a composite swagger
	# this service should follow the same structure as the resources service and clients 
	machinelearning: {version: "2016-05-01-preview", swagger: "webservices"},
	# machinelearning: {
	# 	webservices: {version: "2016-05-01-preview", swagger: "webservices"},
	# 	commitmentplans: {version: "2016-05-01-preview", swagger: "commitmentPlans"}
	# },
	mediaservices: {version: "2015-10-01", swagger: "media"},
	mobileengagement: {version: "2014-12-01", swagger: "mobile-engagement"},
	network: {version: "2016-09-01"},
	notificationhubs: {version: "2016-03-01"},
	powerbiembedded: {version: "2016-01-29"},
	# recoveryservices: {version: "2016-06-01"},
	# recoveryservicesbackup: {version: "2016-06-01"},
	redis: {version: "2016-04-01"},
	resources: {
		features: {version: "2015-12-01"},
		locks: {version: "2016-09-01"},
		policy: {version: "2016-04-01"},
		resources: {version: "2016-09-01"},
		# AutoRest Go generator has a bug and generates an ugly SDK for subscription newest API version
        # https://github.com/Azure/autorest/issues/1477
		# subscriptions: {version: "2016-06-01"}
		subscriptions: {version: "2015-11-01"}
	},
	scheduler: {version: "2016-03-01"},
	search: {version: "2015-02-28"},
	servermanagement: {version: "2016-07-01-preview"},
	servicebus: {version: "2015-08-01"},
	sql: {version: "2015-05-01"},
	storage: {version: "2016-01-01"},
	trafficmanager: {version: "2015-11-01"},
	web: {version: "2015-08-01", swagger: "service"} # composite swagger
}

DATAPLANE_SWAGGERS = {
	# batch: {version: "2016-07-01.3.1", swagger: "BatchService"},
	# insights (composite swagger)
	keyvault: {version: "2015-06-01"},
	# same thing as machine learning, two swaggers but no comp swagger
	# search: {
	# 	searchindex: {version: "2015-02-28"},
	# 	searchservice: {version: "2015-02-28"}
	# },
	# servicefabric: {version: "2016-01-28"}
}

MISPLACED_SWAGGERS = {
	# datalake_analytics: {
	# 	catalog: {version: "2016-06-01-preview"},
	# 	job:  {version: "2016-03-20-preview"}
	# },
	# datalake_store: {
	# 	filesystem: {version: "2015-10-01-preview"}
	# },
}

class Service
	GO_NAMESPACE = "github.com/Azure/azure-sdk-for-go/%s/%s"
	INPUT_PATH   = "%s/azure-rest-api-specs/%s%s/swagger/%s.json"
	OUTPUT_PATH  = "#{ENV['GOPATH']}/src/%s"

	attr :name
	attr :fullname
	attr :namespace
	attr :packages
	attr :task_name
	attr :version
	attr :plane

    attr :spec_root

	attr :input_path
	attr :output_path

	def initialize(service, input_prefix, plane)
		@plane = plane
		@packages = service[:packages]
		@name = @packages.last
		@fullname = @packages.map{|package| package.to_s.gsub(/_/,'-')}.join('/')
		@namespace = sprintf(GO_NAMESPACE, @plane, @fullname)
		@task_name = plane + ':' + @packages.join(':')
		@version = service[:version]        
		if ENV['SWAGGER'] == nil
            spec_root = "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/%s%s/swagger/%s.json"
        else
            spec_root = sprintf("%s/azure-rest-api-specs/%s", ENV['SWAGGER'], "%s%s/swagger/%s.json")
        end
		swagger = service[:swagger] || @name
		@input_path = sprintf(spec_root, input_prefix, [@fullname, @version].join('/'), swagger)
		@output_path = sprintf(OUTPUT_PATH, @namespace)
	end
end

def to_services(m, h, ip, plane, *p)
  if h.keys.include?(:version)
    h[:packages] = p.reverse
    m << Service.new(h, ip, plane)
  else
    h.keys.each do |k|
      to_services(m, h[k], ip, plane, k, *p)
    end
  end
  m
end

services = to_services([], ARM_SWAGGERS, "arm-", "arm")
services = to_services(services, DATAPLANE_SWAGGERS, "", "dataplane")
services = to_services(services, MISPLACED_SWAGGERS, "arm-", "dataplane")

desc "Generate, format, and build all services"
task :default => 'generate:all'

desc "List the known services"
task :services do
	services.each do |service|
		puts "#{service.task_name}"
	end
end

namespace :generate do
	desc "Generate all services"
	task :all do
		services.each {|service| Rake::Task["generate:#{service.task_name}"].execute }
	end

	services.each do |service|
		desc "Generate the #{service.task_name} service"
		task service.task_name.to_sym do
			generate(service)
		end
	end	
end

namespace :go do
	namespace :delete do
		desc "Delete all generated services"
		task :all do
			services.each {|service| delete(service) }
		end

		services.each do |service|
			desc "Delete the #{service.task_name} service"
			task service.task_name.to_sym do
				delete(service)
			end
		end
	end

	namespace :format do
		desc "Format all generated services"
		task :all do
			services.each {|service| format(service) }
		end

		services.each do |service|
			desc "Format the #{service.task_name} service"
			task service.task_name.to_sym do
				format(service)
			end
		end
	end

	namespace :build do
		desc "Build all generated services"
		task :all do
			services.each {|service| build(service) }
		end

		services.each do |service|
			desc "Build the #{service.task_name} service"
			task service.task_name.to_sym do
				build(service)
			end
		end
	end

	namespace :lint do
		desc "Lint all generated services"
		task :all do
			services.each {|service| lint(service) }
		end

		services.each do |service|
			desc "Lint the #{service.task_name} service"
			task service.task_name.to_sym do
				lint(service)
			end
		end
	end

	namespace :vet do
		desc "Vet all generated services"
		task :all do
			services.each {|service| vet(service) }
		end

		services.each do |service|
			desc "Vet the #{service.task_name} service"
			task service.task_name.to_sym do
				vet(service)
			end
		end
	end
end

def generate(service)
	s = "Generating #{service.plane} #{service.fullname}"
	puts "#{s} #{"=" * (80 - s.length)}"
	delete(service)
	s = `#{sprintf(AUTOREST, ENV['AUTOREST'])} -AddCredentials -CodeGenerator Go -Header MICROSOFT_APACHE -Input #{service.input_path} -Namespace #{service.namespace} -OutputDirectory #{service.output_path} -Modeler Swagger -pv #{ENV['SDK_VERSION']}`
	raise "Failed generating #{service.plane} #{service.fullname}.#{service.inspect}" if s =~ /.*FATAL.*/
	puts s

	format(service)
	build(service)
	lint(service)
	vet(service)
end

def delete(service)
	puts "Deleting #{service.plane} #{service.fullname}"
	FileUtils.rmtree(service.output_path)
end

def format(service)
	puts "Formatting #{service.plane} #{service.fullname}"
	s = `gofmt -w #{service.output_path}`
	raise "Formatting #{service.output_path} failed:\n#{s}\n" if $?.exitstatus > 0
end

def build(service)
	puts "Building #{service.plane} #{service.fullname}"
	s = `go build #{service.namespace}`
	raise "Building #{service.namespace} failed:\n#{s}\n" if $?.exitstatus > 0
end

def lint(service)
	puts "Linting #{service.plane} #{service.fullname}"
	s = `#{ENV["GOPATH"]}\\bin\\golint #{service.namespace}`
	print s
end

def vet(service)
	puts "Vetting #{service.plane} #{service.fullname}"
	s = `go vet #{service.namespace}`
	print s
end
