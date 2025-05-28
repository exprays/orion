'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { 
  Search, 
  Key, 
  Trash2, 
  Edit, 
  Clock, 
  Database,
  Hash,
  List,
  Type
} from 'lucide-react';
import { KeyInfo } from '@/types';
import { orionAPI } from '@/lib/api';

export default function KeyManagement() {
  const [keys, setKeys] = useState<KeyInfo[]>([]);
  const [filteredKeys, setFilteredKeys] = useState<KeyInfo[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedKey, setSelectedKey] = useState<KeyInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchKeys();
  }, []);

  useEffect(() => {
    // Filter keys based on search term
    const filtered = keys.filter(key => 
      key.key.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredKeys(filtered);
  }, [keys, searchTerm]);

  const fetchKeys = async () => {
    setIsLoading(true);
    const response = await orionAPI.getKeys('*', 1000);
    
    if (response.success && response.data) {
      setKeys(response.data);
      setError(null);
    } else {
      setError(response.error || 'Failed to fetch keys');
    }
    setIsLoading(false);
  };

  const handleKeySelect = async (keyName: string) => {
    const response = await orionAPI.getKey(keyName);
    if (response.success && response.data) {
      setSelectedKey(response.data);
    }
  };

  const handleKeyDelete = async (keyName: string) => {
    if (!confirm(`Are you sure you want to delete key "${keyName}"?`)) {
      return;
    }

    const response = await orionAPI.deleteKey(keyName);
    if (response.success) {
      setKeys(prev => prev.filter(key => key.key !== keyName));
      if (selectedKey?.key === keyName) {
        setSelectedKey(null);
      }
    } else {
      alert('Failed to delete key: ' + response.error);
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'string': return <Type className="h-4 w-4" />;
      case 'set': return <List className="h-4 w-4" />;
      case 'hash': return <Hash className="h-4 w-4" />;
      default: return <Database className="h-4 w-4" />;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'string': return 'bg-blue-100 text-blue-800';
      case 'set': return 'bg-green-100 text-green-800';
      case 'hash': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const formatTTL = (ttl: number): string => {
    if (ttl === -1) return 'No expiry';
    if (ttl === -2) return 'Expired';
    
    const hours = Math.floor(ttl / 3600);
    const minutes = Math.floor((ttl % 3600) / 60);
    const seconds = ttl % 60;
    
    if (hours > 0) return `${hours}h ${minutes}m ${seconds}s`;
    if (minutes > 0) return `${minutes}m ${seconds}s`;
    return `${seconds}s`;
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="h-8 bg-gray-200 rounded w-1/4 animate-pulse"></div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="animate-pulse">
            <CardContent className="pt-6">
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="h-12 bg-gray-200 rounded"></div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Key Management</h2>
        <Button onClick={fetchKeys} variant="outline" size="sm">
          Refresh
        </Button>
      </div>

      {/* Search */}
      <Card>
        <CardContent className="pt-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search keys..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>
          <div className="mt-2 text-sm text-gray-500">
            Showing {filteredKeys.length} of {keys.length} keys
          </div>
        </CardContent>
      </Card>

      {error && (
        <Card className="border-red-200 bg-red-50">
          <CardContent className="pt-6">
            <p className="text-red-600">{error}</p>
          </CardContent>
        </Card>
      )}

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Keys List */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Key className="h-5 w-5" />
              <span>Keys ({filteredKeys.length})</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="max-h-96 overflow-y-auto">
            <div className="space-y-2">
              {filteredKeys.map((key) => (
                <div
                  key={key.key}
                  className={`p-3 border rounded-lg cursor-pointer transition-colors hover:bg-gray-50 ${
                    selectedKey?.key === key.key ? 'border-blue-500 bg-blue-50' : ''
                  }`}
                  onClick={() => handleKeySelect(key.key)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2 flex-1 min-w-0">
                      {getTypeIcon(key.type)}
                      <span className="font-mono text-sm truncate">{key.key}</span>
                      <Badge className={`text-xs ${getTypeColor(key.type)}`}>
                        {key.type}
                      </Badge>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleKeyDelete(key.key);
                      }}
                      className="text-red-500 hover:text-red-700"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                  <div className="mt-1 flex items-center space-x-4 text-xs text-gray-500">
                    <span>Size: {key.size} bytes</span>
                    <span className="flex items-center space-x-1">
                      <Clock className="h-3 w-3" />
                      <span>TTL: {formatTTL(key.ttl)}</span>
                    </span>
                  </div>
                </div>
              ))}
              
              {filteredKeys.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  {searchTerm ? 'No keys match your search' : 'No keys found'}
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Key Details */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Edit className="h-5 w-5" />
              <span>Key Details</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            {selectedKey ? (
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium text-gray-600">Key Name</label>
                  <p className="font-mono text-sm bg-gray-100 p-2 rounded">{selectedKey.key}</p>
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium text-gray-600">Type</label>
                    <div className="flex items-center space-x-2 mt-1">
                      {getTypeIcon(selectedKey.type)}
                      <Badge className={getTypeColor(selectedKey.type)}>
                        {selectedKey.type}
                      </Badge>
                    </div>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">TTL</label>
                    <p className="text-sm mt-1">{formatTTL(selectedKey.ttl)}</p>
                  </div>
                </div>

                <div>
                  <label className="text-sm font-medium text-gray-600">Value</label>
                  <div className="mt-1">
                    {selectedKey.type === 'string' && (
                      <textarea
                        className="w-full p-2 border rounded font-mono text-sm"
                        rows={4}
                        value={selectedKey.value || ''}
                        readOnly
                      />
                    )}
                    
                    {selectedKey.type === 'set' && (
                      <div className="space-y-1">
                        {selectedKey.members?.map((member, index) => (
                          <div key={index} className="font-mono text-sm bg-gray-100 p-2 rounded">
                            {member}
                          </div>
                        ))}
                      </div>
                    )}
                    
                    {selectedKey.type === 'hash' && (
                      <div className="space-y-1">
                        {Object.entries(selectedKey.fields || {}).map(([field, value]) => (
                          <div key={field} className="grid grid-cols-2 gap-2">
                            <div className="font-mono text-sm bg-gray-100 p-2 rounded">{field}</div>
                            <div className="font-mono text-sm bg-gray-100 p-2 rounded">{value}</div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                Select a key to view its details
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}